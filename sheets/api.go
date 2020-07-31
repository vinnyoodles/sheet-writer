package sheets

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"strings"
)

type Config struct {
	DatabaseID string `json:"db_id"`
}

var RANGE_DELIMITER = ':'
var DriveService *drive.Service
var SheetsService *sheets.Service
var SheetDatabaseID string

func init() {
	DriveService, SheetsService = CreateServices()
	byt, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed reading config file: %v", err)
		return
	}
	config := Config{}
	if err := json.Unmarshal(byt, &config); err != nil {
		log.Fatalf("Failed to unmarshall config file: %v", err)
		return
	}
	SheetDatabaseID = config.DatabaseID
}

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

func Fetch() []string {
	output := []string{}
	query := "mimeType='application/vnd.google-apps.spreadsheet'"
	response, err := DriveService.Files.List().Q(query).Do()
	if err != nil {
		return output
	}
	for i := 0; i < len(response.Files); i++ {
		output = append(output, response.Files[i].Id)
	}
	return output
}

func FetchID(name string) string {
	headers := Headers(SheetDatabaseID)
	for i := 0; i < len(headers); i++ {
		if headers[i] != name {
			continue
		}
		i += 97
		target := fmt.Sprintf("%c2", rune(i))
		result := Read(SheetDatabaseID, target)
		if len(result) == 0 {
			continue
		}
		return result[0]
	}
	return ""
}

func Append(sheetID string, values [][]interface{}) bool {
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
	_, err := SheetsService.Spreadsheets.Values.Append(sheetID, sheetRange, body).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
		return false
	}
	return true
}

func Headers(id string) []string {
	return Read(id, "A1:1")
}

func Read(id string, sheetRange string) []string {
	output := []string{}
	response, err := SheetsService.Spreadsheets.Values.Get(id, sheetRange).Do()
	if err != nil || len(response.Values) == 0 {
		return output
	}
	// TODO: handle multi dimensional arrays
	row := response.Values[0]
	for _, cell := range row {
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
