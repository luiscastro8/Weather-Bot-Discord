package discord

import (
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/weather"
	"Weather-Bot-Discord/weather/forecast"
	"Weather-Bot-Discord/weather/points"
	"Weather-Bot-Discord/weather/zip"
	"github.com/bwmarrin/discordgo"
)

func zipHandler(s *discordgo.Session, i *discordgo.InteractionCreate, zipCode string) {
	if !isValidZip(zipCode) {
		err := sendSlashCommandResponse(s, i, "Error: Zip code must be exactly 5 digits long")
		if err != nil {
			mylogger.Errorln("could not send slash command message:", err)
		}
		return
	}

	forecastUrl, ok := weather.GetUrlFromCache(zipCode)
	if !ok {
		lat, long, err := zip.GetCoordsFromZip(zipCode)
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		weather.AcquireLockForCaching()
		forecastUrl, err = points.GetForecastURLFromCoords(lat, long)
		if err != nil {
			weather.ReleaseLockForCaching()
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}
		weather.WriteToCache(zipCode, forecastUrl)
	}

	forecastMessage, err := forecast.GetForecastFromURL(forecastUrl, "")
	if err != nil {
		mylogger.Errorln(err)
		err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
		if err != nil {
			mylogger.Errorln("could not send slash command message:", err)
		}
		return
	}

	err = sendSlashCommandResponse(s, i, forecastMessage)
	if err != nil {
		mylogger.Errorln("could not send slash command message:", err)
	} else {
		mylogger.Println("sent weather forecast for zip", zipCode)
	}
}
