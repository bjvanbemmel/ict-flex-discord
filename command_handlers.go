package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var (
    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        "register-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            options := i.ApplicationCommandData().Options

            optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
            for _, opt := range options {
                optionMap[opt.Name] = opt
            }

            TargetRole = optionMap["role"].RoleValue(nil, "")
            log.Info("Role has been configured.")

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
            log.Info("Channel has been configured.")

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Flags: discordgo.MessageFlagsEphemeral,
                    Content: fmt.Sprintf("Registered %v as the target channel.", TargetChannel.Mention()),
                },
            })
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
    }
)
