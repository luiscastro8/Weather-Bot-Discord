package discord

import (
	"Weather-Bot-Discord/api/weather/forecast"
	"Weather-Bot-Discord/api/weather/points"
	"Weather-Bot-Discord/mylogger"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func coordinatesHandler(s *discordgo.Session, i *discordgo.InteractionCreate, lat, long float64) {
	forecastUrl, err := points.GetForecastURLFromCoords(fmt.Sprintf("%.4f", lat), fmt.Sprintf("%.4f", long))
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	forecastMessage, err := forecast.GetForecastFromURL(forecastUrl, "")
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	err = sendSlashCommandResponse(s, i, forecastMessage)
	if err != nil {
		mylogger.Errorln("could not send slash command message:", err)
	} else {
		mylogger.Println("sent weather forecast for coords", lat, long)
	}
}
