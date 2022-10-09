package discord

import (
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/weather"
	"Weather-Bot-Discord/weather/forecast"
	"Weather-Bot-Discord/weather/points"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var WeatherCommand = &discordgo.ApplicationCommand{
	Name:        "weather",
	Description: "get the weather",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "zip",
			Description: "get weather by zip code",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "zip-code",
					Description: "5 digit zip code",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
			},
		},
	},
}

func WeatherHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	subCommand := i.ApplicationCommandData().Options[0]
	if subCommand.Name == "zip" {
		zipCode := strconv.Itoa(int(subCommand.Options[0].IntValue()))
		if len(zipCode) != 5 {
			sendSlashCommandResponse(s, i, "Error: Zip code must have exactly 5 digits")
			return
		}
		for _, c := range zipCode {
			if c < '0' || c > '9' {
				sendSlashCommandResponse(s, i, "Error: Zip code must contain all numbers")
				return
			}
		}

		lat, long, err := weather.GetCoordsFromZip(zipCode)
		if err != nil {
			mylogger.Get().Errorln(err)
			sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			return
		}

		url, err := points.GetForecastURLFromCoords(lat, long)
		if err != nil {
			mylogger.Get().Errorln(err)
			sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			return
		}

		forecastMessage, err := forecast.GetForecastFromURL(url)
		if err != nil {
			mylogger.Get().Errorln(err)
			sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			return
		}

		sendSlashCommandResponse(s, i, forecastMessage)
	}
}

func sendSlashCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		mylogger.Get().Errorln("could not send slash command message:", err)
	}
}
