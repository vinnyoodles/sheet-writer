package sheets

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"strings"
)

func Scopes() []string {
	return []string{
		"https://www.googleapis.com/auth/drive.readonly",
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	}
}

func CreateServices() (*drive.Service, *sheets.Service) {
	data, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Failed reading credentials file: %v", err)
		return nil, nil
	}
	scopes := Scopes()
	config, err := google.JWTConfigFromJSON(data, strings.Join(scopes, " "))
	if err != nil {
		log.Fatalf("Failed creating config: %v", err)
		return nil, nil
	}
	client := config.Client(oauth2.NoContext)
	sheetsService, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Failed creating sheets client: %v", err)
		return nil, nil
	}
	driveService, err := drive.New(client)
	if err != nil {
		log.Fatalf("Failed creating drive client: %v", err)
		return nil, nil
	}
	return driveService, sheetsService
}

func Fetch(drv *drive.Service) []string {
	output := []string{}
	query := "mimeType='application/vnd.google-apps.spreadsheet'"
	response, err := drv.Files.List().Q(query).Do()
 	if err != nil {
		return output
 	}
 	for i := 0; i < len(response.Files); i++ {
		output = append(output, response.Files[i].Id)
 	}
	return output
}

func Create() bool {
	return false
}

func Write() bool {
	return false
}

func Append() bool {
	return false
}

func RowCount() int {
	return 0
}

func Headers(srv *sheets.Service, id string) []string {
	output := []string{}
	sheetRange := "A1:1"
	response, err := srv.Spreadsheets.Values.Get(id, sheetRange).Do()
	if err != nil || len(response.Values) == 0 {
		return output
	}
	headerRow := response.Values[0]
	for _, cell := range headerRow {
		output = append(output, cell.(string))
	}
	return output
}

