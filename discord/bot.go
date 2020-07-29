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
	var data map[string]interface{}
	byt, err := ioutil.ReadFile("bot.json")
	if err != nil {
		log.Fatalf("Failed reading bot file: %v", err)
		return nil
	}
	if err := json.Unmarshal(byt, &data); err != nil {
		log.Fatalf("Failed reading bot file: %v", err)
		return nil
	}
	token := data["token"]
	session, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatalf("Failed to create discord session")
		return nil
	}
	return session
}