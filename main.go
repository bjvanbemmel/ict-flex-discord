package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var (
    Token *string = flag.String("t", "", "Bot token")
    GUID *string = flag.String("g", "", "Guild ID")

    TargetRole *discordgo.Role
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

var (
    commands = []*discordgo.ApplicationCommand{
        {
            Name: "marco",
            Description: "You say Marco, I say...",
        },
        {
            Name: "link",
            Description: "Shows you a link thingamaling",
        },
        {
            Name: "register-channel",
            Description: "Register a channel that you would like the announcements to be posted in",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type: discordgo.ApplicationCommandOptionChannel,
                    Name: "channel",
                    Description: "The target channel",
                    ChannelTypes: []discordgo.ChannelType{
                        discordgo.ChannelTypeGuildText,
                        discordgo.ChannelTypeGuildNews,
                    },
                    Required: true,
                },
            },
        },
        {
            Name: "register-role",
            Description: "Register the role you want to be pinged to notify users about new announcements",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type: discordgo.ApplicationCommandOptionRole,
                    Name: "role",
                    Description: "The target role",
                    Required: true,
                },
            },
        },
        {
            Name: "embed",
            Description: "Sends an embed... I think..?",
        },
        {
            Name: "inspect",
            Description: "Shows the current configuration",
        },
    }

    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        "marco": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: "Polo!",
                },
            })
        },
        "link": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: "@everyone Title\nDescription",
                    Flags: discordgo.MessageFlagsEphemeral,
                    Components: []discordgo.MessageComponent{
                        discordgo.ActionsRow{
                            Components: []discordgo.MessageComponent{
                                discordgo.Button{
                                    Label: "Article",
                                    Style: discordgo.LinkButton,
                                    URL: "https://github.com/bjvanbemmel/rtouch/releases/latest",
                                },
                            },
                        },
                    },
                },
            })
            if err != nil {
                log.Fatal(err.Error())
            }
        },
        "register-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            options := i.ApplicationCommandData().Options

            optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
            for _, opt := range options {
                optionMap[opt.Name] = opt
            }

            TargetRole = optionMap["role"].RoleValue(nil, "")

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Flags: discordgo.MessageFlagsEphemeral,
                    Content: fmt.Sprintf("Registered %v as the target role.", TargetRole.Mention()),
                },
            })
        },
        "register-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            options := i.ApplicationCommandData().Options

            optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
            for _, opt := range options {
                optionMap[opt.Name] = opt
            }

            TargetChannel = optionMap["channel"].ChannelValue(nil)

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Flags: discordgo.MessageFlagsEphemeral,
                    Content: fmt.Sprintf("Registered %v as the target channel.", TargetChannel.Mention()),
                },
            })
        },
        "embed": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            t, err := time.Parse("02-01-2006", "19-04-2023")
            if err != nil {
                log.Fatal(err.Error())
            }

            embed := discordgo.MessageEmbed{
                Title: "Nieuw onderzoek stagediscriminatie",
                URL: "https://ict-flex.nl/nieuw-onderzoek-stagediscriminatie",
                Type: discordgo.EmbedTypeArticle,
                Description: "In het verlengde van het pas ondertekende Stagepact mbo onderzoekt het Verwey-Jonker Instituut behoeften en wensen van mbo-studenten rondom het melden van stagediscriminatie. Doel: nagaan óf die behoefte onder studenten leeft. En zo ja, hoe zij dat het liefst doen, bij wie en onder welke voorwaarden, om vervolgens beleidsmakers gericht te adviseren. Het onderzoeksteam komt graag in contact met studenten die willen meepraten. Wil je meedoen, klik dan hieronder voor het aanmeldformulier.",
                Timestamp: t.Format(time.RFC3339),
                Color: 0xEF7E05,
                Author: &discordgo.MessageEmbedAuthor{
                    URL: "https://ict-flex.nl/author/adalmolen/",
                    Name: "André Dalmolen",
                },
            }

            msg := discordgo.MessageSend{
                Content: fmt.Sprintf("%v New announcement!", TargetRole.Mention()),
                Embeds: []*discordgo.MessageEmbed{
                    &embed,
                },
                Components: []discordgo.MessageComponent{
                    discordgo.ActionsRow{
                        Components: []discordgo.MessageComponent{
                            discordgo.Button{
                                Label: "Aanmelden",
                                Style: discordgo.LinkButton,
                                URL: "https://ict-flex.nl/nieuw-onderzoek-stagediscriminatie/",
                            },
                        },
                    },
                },
            }

            s.ChannelMessageSendComplex(i.ChannelID, &msg)
        },
        "inspect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            var tr string = "<nil>"
            var tc string = "<nil>"

            if TargetRole != nil {
                tr = TargetRole.Mention()
            }

            if TargetChannel != nil {
                tc = TargetChannel.Mention()
            }

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Flags: discordgo.MessageFlagsEphemeral,
                    Embeds: []*discordgo.MessageEmbed{
                        {
                            Type: discordgo.EmbedTypeRich,
                            Title: "Configuration",
                            Description: "This is the current configuration:",
                            Color: 0xEF7E05,
                            Fields: []*discordgo.MessageEmbedField{
                                {
                                    Name: "Target channel:",
                                    Value: tc,
                                },
                                {
                                    Name: "Target role:",
                                    Value: tr,
                                },
                            },
                        },
                    },
                },
            })
        },
    }
)

func init() {
    session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if handle, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
            handle(s, i)
        }
    })
}

func main() {
    session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
        log.Print("Logged in as: %s#%v", s.State.User.Username, s.State.User.Discriminator)
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
