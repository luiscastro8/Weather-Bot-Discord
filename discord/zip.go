package discord

import (
	"Weather-Bot-Discord/api/weather/forecast"
	"Weather-Bot-Discord/api/weather/points"
	"Weather-Bot-Discord/mylogger"
	zip2 "Weather-Bot-Discord/zip"
	"github.com/bwmarrin/discordgo"
)

func zipHandler(s *discordgo.Session, i *discordgo.InteractionCreate, zipCode string, hourly bool) {
	if !isValidZip(zipCode) {
		sendSlashCommandResponseAndLogError(s, i, "Error: Zip code must be exactly 5 digits long")
		return
	}

	forecastUrl, err := getForcecastURL(zipCode, hourly)
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	forecastMessage, err := forecast.GetForecastFromURL(forecastUrl, "", hourly)
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	err = sendSlashCommandResponse(s, i, forecastMessage)
	if err != nil {
		mylogger.Errorln("could not send slash command message:", err)
	} else {
		mylogger.Println("sent weather forecast for zip", zipCode)
	}
}

func isValidZip(zipCode string) bool {
	if len(zipCode) != 5 {
		return false
	}
	for _, c := range zipCode {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func getForcecastURL(zipCode string, hourly bool) (string, error) {
	forecastUrl := ""
	ok := false

	if hourly {
		forecastUrl, ok = zip2.GetHourlyUrlFromCache(zipCode)
		if !ok {
			lat, long, err := zip2.GetCoordsFromZip(zipCode)
			if err != nil {
				return "", err
			}

			zip2.AcquireHourlyLockForCaching()
			forecastUrl, err = points.GetHourlyForecastURLFromCoords(lat, long)
			if err != nil {
				zip2.ReleaseHourlyLockForCaching()
				return "", err
			}
			zip2.WriteToHourlyCache(zipCode, forecastUrl)
		}
	}

	forecastUrl, ok = zip2.GetDailyUrlFromCache(zipCode)
	if !ok {
		lat, long, err := zip2.GetCoordsFromZip(zipCode)
		if err != nil {
			return "", err
		}

		zip2.AcquireDailyLockForCaching()
		forecastUrl, err = points.GetDailyForecastURLFromCoords(lat, long)
		if err != nil {
			zip2.ReleaseDailyLockForCaching()
			return "", err
		}
		zip2.WriteToDailyCache(zipCode, forecastUrl)
	}

	return forecastUrl, nil
}
