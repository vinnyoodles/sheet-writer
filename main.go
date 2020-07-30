package main

import (
	"fmt"
	"github.com/sheetwriter/discord"
	"github.com/sheetwriter/sheets"
)

func main() {
	drv, srv := sheets.CreateServices()
	sheetIDs := sheets.Fetch(drv)
	fmt.Println(sheetIDs)
	for i := 0; i < len(sheetIDs); i++ {
		headers := sheets.Headers(srv, sheetIDs[i])
		fmt.Println(headers)
		stringValues := [2][2]string{
			{"abc", "123"},
			{"def", "456"},
		}
		values := make([][]interface{}, len(stringValues))
		for idx, arr := range stringValues {
			values[idx] = make([]interface{}, len(arr))
			for idx2, s := range arr {
				values[idx][idx2] = s
			}
		}
		// result := sheets.Append(srv, sheetIDs[i], values)
		// fmt.Println(result)
	}

	session := discord.CreateSession()
	discord.Run(session)
}
