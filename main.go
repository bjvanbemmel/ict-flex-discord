package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bjvanbemmel/ict-flex-discord/scraper"
	"github.com/bjvanbemmel/ict-flex-discord/storage"
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var (
	Token string = os.Getenv("USER_TOKEN")
	GUID  string = os.Getenv("GUILD_ID")

	Settings *storage.Settings = &storage.Settings{}
)

var session *discordgo.Session

func init() {
	flag.Parse()

	var err error
	session, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handle, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handle(s, i)
		}
	})
}

func main() {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %s#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := session.Open()
	if err != nil {
		log.Fatal(err.Error())
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, GUID, v)
		if err != nil {
			log.Fatal(err.Error())
		}

		registeredCommands[i] = cmd
	}

	err = session.UpdateGameStatus(0, "Please contribute to the project!")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer session.Close()

	var scraper scraper.Scraper = scraper.Scraper{}
	go func() {
		for {
			time.Sleep(time.Second * 30)

			role, err := Settings.Role()
			if err != nil {
				log.Errorf("Something went wrong while fetching role: `%s`", err.Error())
			}

			channel, err := Settings.Channel()
			if err != nil {
				log.Errorf("Something went wrong while fetching channel: `%s`", err.Error())
			}

			if role == nil || channel == nil {
				log.Infof("Values nil. Skipping... `%v`", Settings)
				continue
			}

			embeds, err := scraper.Start()
			if err != nil {
				log.Fatal(err.Error())

				return
			}

			raw, _ := json.Marshal(embeds)
			fmt.Println(string(raw))

			if len(embeds) == 0 {
				continue
			}

			for _, embed := range embeds {
				msg := discordgo.MessageSend{
					Content: fmt.Sprintf("%v New announcement!", Settings.TargetRole.Mention()),
					Embeds: []*discordgo.MessageEmbed{
						&embed,
					},
				}

				log.Info("Sending embed to TargetChannel...")
				_, err := session.ChannelMessageSendComplex(Settings.TargetChannel.ID, &msg)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Print("Press CTRL+C to exit")
	<-stop

	log.Info("Removing commands...")

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, GUID, v.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
