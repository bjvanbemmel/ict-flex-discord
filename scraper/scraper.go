package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bjvanbemmel/ict-flex-rss/types"
	"github.com/bwmarrin/discordgo"
)

type Scraper struct {
	Previous types.Feed
}

func (s *Scraper) Start() ([]discordgo.MessageEmbed, error) {
	var client *http.Client = &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://rss.bjvanbemmel.nl/ict-flex", nil)
	req.Header.Set("User-Agent", "Beep boop. I am the ICT-Flex-Discord bot. Boop bap.")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	re, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed types.Feed

	header := len("<rss xmlns:atom\"http://www.w3.org/2005/Atom\" version=\"2.0\">")
	footer := len("</rss>")

	// Ignore the header and footer
	err = xml.Unmarshal(re[header:len(re)-footer], &feed)
	if err != nil {
		return nil, err
	}

	if s.Previous.Title == "" {
		s.Previous = feed

		return nil, nil
	}

	var newArticles []*types.Article = []*types.Article{}
	for _, art := range feed.Articles {

		var found bool = false
		for _, oldArt := range s.Previous.Articles {
			if oldArt.Guid == art.Guid {
				found = true

				break
			}
		}

		if found {
			continue
		}

		fmt.Printf("[NEW] [%s] GUID: `%v` TITLE: `%s`\n", time.Now().String(), art.Guid, art.Title)
		fmt.Println("OLD ARTICLES:")
		for _, art := range s.Previous.Articles {
			fmt.Printf("GUID: `%v` TITLE: `%s`\n", art.Guid, art.Title)
		}
		fmt.Println("END OF OLD ARTICLES")
		fmt.Println("START OF RSS FEED")
		fmt.Println(string(re))
		fmt.Println("END OF RSS FEED")

		newArticles = append(newArticles, art)
	}

	embeds := s.CreateEmbeds(newArticles)
	s.Previous = feed

	return embeds, nil
}

func (s Scraper) CreateEmbeds(articles []*types.Article) []discordgo.MessageEmbed {
	var embeds []discordgo.MessageEmbed = []discordgo.MessageEmbed{}

	for _, art := range articles {
		embed := discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeArticle,
			Title:       art.Title,
			URL:         art.Link,
			Description: art.Description,
			Timestamp:   art.CreatedAt.Format(time.RFC3339),
			Author: &discordgo.MessageEmbedAuthor{
				Name: art.Author,
			},
		}

		embeds = append(embeds, embed)
	}

	return embeds
}
