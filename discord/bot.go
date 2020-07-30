package discord

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type BotCredentials struct {
	Token string `json:"token"`
}

func CreateSession() *discordgo.Session {
	byt, err := ioutil.ReadFile("bot.json")
	if err != nil {
		log.Fatalf("Failed reading bot file: %v", err)
		return nil
	}
	creds := BotCredentials{}
	if err := json.Unmarshal(byt, &creds); err != nil {
		log.Fatalf("Failed reading bot file: %v", err)
		return nil
	}
	session, err := discordgo.New(fmt.Sprintf("Bot %s", creds.Token))
	if err != nil {
		log.Fatalf("Failed to create discord session: %v", err)
		return nil
	}
	session.AddHandler(messageListen)
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	return session
}

func Run(session *discordgo.Session) {
	err := session.Open()
	if err != nil {
		log.Fatalf("Failed to open discord socket: %v", err)
		return
	}
	// Listen to syscalls to stop running
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Wait until one of the above signals are sent
	<-channel
	session.Close()
}

func messageListen(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		// ignore bot messages
		return
	}
	content := message.Content
	fmt.Println(content)
}
