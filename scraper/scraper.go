package scraper

import (
	"encoding/xml"
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
	resp, err := http.Get("https://rss.bjvanbemmel.nl/ict-flex")
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
				URL:  art.Author.Profile,
				Name: art.Author.Name,
			},
		}

		embeds = append(embeds, embed)
	}

	return embeds
}
