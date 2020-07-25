package main

import (
	"fmt"
	"github.com/sheetwriter/sheets"
)

func main() {
	drv, srv := sheets.CreateServices()
	sheetIDs := sheets.Fetch(drv)
	for i := 0; i < len(sheetIDs); i ++ {
		headers := sheets.Headers(srv, sheetIDs[i])
		fmt.Println(headers)
	}
}
