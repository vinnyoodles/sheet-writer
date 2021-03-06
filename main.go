package main

import (
	"fmt"
	"github.com/sheetwriter/discord"
	"github.com/sheetwriter/sheets"
	"time"
)

var DateTimeFormat = "1/2/2006 03:04 PM"

func main() {
	session := discord.CreateSession()
	discord.AddHandler("AppendSheet", AppendSheet)
	fmt.Println("Starting bot")
	discord.Run(session)
}

func AppendSheet(args []string) {
	if len(args) < 1 {
		return
	}
	name := args[0]
	sheetID := sheets.FetchID(name)
	if len(sheetID) == 0 {
		return
	}
	values := make([]interface{}, len(args))

	location, err := time.LoadLocation("US/Pacific")
	t := time.Now()
	if err == nil {
		t = t.In(location)
	}
	for idx, val := range args {
		if idx == 0 {
			values[idx] = t.Format(DateTimeFormat)
		} else {
			values[idx] = val
		}
	}
	wrapper := make([][]interface{}, 1)
	wrapper[0] = values
	fmt.Println("Wrote:", wrapper)
	sheets.Append(sheetID, wrapper)
}
