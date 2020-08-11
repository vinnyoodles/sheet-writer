package discord

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type BotCredentials struct {
	Token string `json:"discord_token"`
}

type Commands struct {
	Commands []BotCommand `json:"commands"`
}

type BotCommand struct {
	Function string `json:"function"`
	Prefix   string `json:"prefix"`
	Name     string `json:"name"`
}

type CommandHandler func([]string)

var CommandMap map[string]BotCommand
var FunctionMap map[string]CommandHandler

func init() {
	byt, err := ioutil.ReadFile("commands.json")
	if err != nil {
		log.Fatalf("Failed reading commands file: %v", err)
		return
	}
	CommandMap = make(map[string]BotCommand)
	commands := Commands{}
	json.Unmarshal(byt, &commands)
	for i := 0; i < len(commands.Commands); i++ {
		cmd := commands.Commands[i]
		CommandMap[cmd.Prefix] = cmd
	}

	FunctionMap = make(map[string]CommandHandler)
}

func CreateSession() *discordgo.Session {
	var token string
	val, ok := os.LookupEnv("DISCORD_TOKEN")
	if ok {
		token = val
	} else {
		byt, err := ioutil.ReadFile("config.json")
		if err != nil {
			log.Fatalf("Failed reading config file: %v", err)
			return nil
		}
		creds := BotCredentials{}
		if err := json.Unmarshal(byt, &creds); err != nil {
			log.Fatalf("Failed to unmarshal config file: %v", err)
			return nil
		}
		token = creds.Token
	}

	session, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatalf("Failed to create discord session: %v", err)
		return nil
	}
	session.AddHandler(messageListen)
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	return session
}

func AddHandler(name string, fn CommandHandler) {
	FunctionMap[name] = fn
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
		// Ignore bot messages
		return
	}
	content := message.Content
	tokens := strings.Fields(content)
	if len(tokens) == 0 {
		return
	}

	cmd, ok := CommandMap[tokens[0]]
	if !ok {
		return
	}
	args := tokens[1:]
	if len(cmd.Name) > 0 {
		args = append([]string{cmd.Name}, args...)

	}
	FunctionMap[cmd.Function](args)
}
