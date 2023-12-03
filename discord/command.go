package discord

import "github.com/bwmarrin/discordgo"

var WeatherCommand = &discordgo.ApplicationCommand{
	Name:        "weather",
	Description: "get the weather",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "daily",
			Description: "daily weather",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "zip",
					Description: "get weather by zip code",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "zip-code",
							Description: "5 digit zip code",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "coordinates",
					Description: "get weather using latitude and longitude",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "latitude",
							Description: "latitude",
							Type:        discordgo.ApplicationCommandOptionNumber,
							Required:    true,
						},
						{
							Name:        "longitude",
							Description: "longitude",
							Type:        discordgo.ApplicationCommandOptionNumber,
							Required:    true,
						},
					},
				},
				{
					Name:        "address",
					Description: "get weather using an address",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "address",
							Description: "address",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "hourly",
			Description: "hourly weather",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "zip",
					Description: "get weather by zip code",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "zip-code",
							Description: "5 digit zip code",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "coordinates",
					Description: "get weather using latitude and longitude",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "latitude",
							Description: "latitude",
							Type:        discordgo.ApplicationCommandOptionNumber,
							Required:    true,
						},
						{
							Name:        "longitude",
							Description: "longitude",
							Type:        discordgo.ApplicationCommandOptionNumber,
							Required:    true,
						},
					},
				},
				{
					Name:        "address",
					Description: "get weather using an address",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "address",
							Description: "address",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	},
}
