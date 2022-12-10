package main

import (
	"Weather-Bot-Discord/discord"
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/token"
	"Weather-Bot-Discord/zip"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mylogger.Init(os.Stdout, os.Stderr)

	botToken, err := token.GetToken()
	if err != nil {
		mylogger.Fatalln("Unable to get bot token:", err)
	}

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		mylogger.Fatalln("Unable to create discord session:", err)
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "weather" {
			discord.WeatherHandler(s, i)
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		mylogger.Fatalln("Unable to open discord connection:", err)
	}

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", discord.WeatherCommand)
	if err != nil {
		mylogger.Fatalln("Cannot create command", err)
	}

	err = zip.OpenZipFile("zip-codes.csv")
	if err != nil {
		mylogger.Fatalln("Unable to cache zip codes and coordinates:", err)
	}

	mylogger.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = dg.Close()
	if err != nil {
		mylogger.Fatalln("Unable to close discord connection:", err)
	}
	mylogger.Println("Closed discord connection")
}
