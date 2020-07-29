 package discord

import (
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
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
		log.Fatalf("Failed to create discord session")
		return nil
	}
	return session
}