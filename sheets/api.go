package sheets

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"strings"
)

var RANGE_DELIMITER = ':'

func Scopes() []string {
	return []string{
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/spreadsheets",
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

func Append(srv *sheets.Service, sheetID string, values [][]interface{}) bool {
	ctx := context.Background()
	valueInputOption := "USER_ENTERED"
	insertDataOption := "INSERT_ROWS"
	maxLength := 0
	for i := 0; i < len(values); i++ {
		if len(values[i]) > maxLength {
			maxLength = len(values[i])
		}
	}
	sheetRange := _getRange(0 /*start index*/, maxLength)
	body := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          sheetRange,
		Values:         values,
	}
	_, err := srv.Spreadsheets.Values.Append(sheetID, sheetRange, body).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
		return false
	}
	return true
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

func _getRange(startIndex int, length int) string {
	startIndex += 97
	start := rune(startIndex)
	end := rune(startIndex + length)
	return string([]rune{start, rune(RANGE_DELIMITER), end})
}
