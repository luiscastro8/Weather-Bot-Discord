package discord

import (
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/weather"
	"Weather-Bot-Discord/weather/forecast"
	"Weather-Bot-Discord/weather/points"
	"Weather-Bot-Discord/weather/zip"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
		zipCode := fmt.Sprintf("%05d", subCommand.Options[0].IntValue())
		if len(zipCode) != 5 {
			sendSlashCommandResponse(s, i, "Error: Zip code must have exactly 5 digits")
			return
		}

		forecastUrl, ok := weather.GetUrlFromCache(zipCode)
		if !ok {
			lat, long, err := zip.GetCoordsFromZip(zipCode)
			if err != nil {
				mylogger.Errorln(err)
				sendSlashCommandResponse(s, i, "There was an error getting the forecast")
				return
			}

			weather.AcquireLockForCaching()
			forecastUrl, err = points.GetForecastURLFromCoords(lat, long)
			if err != nil {
				weather.ReleaseLockForCaching()
				mylogger.Errorln(err)
				sendSlashCommandResponse(s, i, "There was an error getting the forecast")
				return
			}
			weather.WriteToCache(zipCode, forecastUrl)
		}

		forecastMessage, err := forecast.GetForecastFromURL(forecastUrl)
		if err != nil {
			mylogger.Errorln(err)
			sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			return
		}

		mylogger.Println("sent weather forecast for zip", zipCode)
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
		mylogger.Errorln("could not send slash command message:", err)
	}
}
