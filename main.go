package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bjvanbemmel/ict-flex-discord/scraper"
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var (
	Token *string = flag.String("t", "", "Bot token")
	GUID  *string = flag.String("g", "", "Guild ID")

	TargetRole    *discordgo.Role
	TargetChannel *discordgo.Channel
)

var session *discordgo.Session

func init() {
	flag.Parse()

	var err error
	session, err = discordgo.New("Bot " + *Token)
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
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, *GUID, v)
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

			if TargetRole == nil || TargetChannel == nil {
				log.Info("Values nil. Skipping...")
				continue
			}

			embeds, err := scraper.Start()
			if err != nil {
				log.Fatal(err.Error())

				return
			}

			if len(embeds) == 0 {
				continue
			}

			for _, embed := range embeds {
				msg := discordgo.MessageSend{
					Content: fmt.Sprintf("%v New announcement!", TargetRole.Mention()),
					Embeds: []*discordgo.MessageEmbed{
						&embed,
					},
				}

				log.Info("Sending embed to TargetChannel...")
				_, err := session.ChannelMessageSendComplex(TargetChannel.ID, &msg)
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
		err := session.ApplicationCommandDelete(session.State.User.ID, *GUID, v.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
