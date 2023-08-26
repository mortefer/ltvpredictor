package parsers

import (
	"encoding/json"
	"io"
	"log"
	"mortefer/ltvpredictor_light/internal/entities"
	"mortefer/ltvpredictor_light/internal/errors"
	"os"
)

/**
 *	JsonParser, implements InputFileParser interface
 */
type JsonParser struct {
}

/**
 *	Parse supplied JSON file, while checking for validity of JSON data
 *  Receives fileName string
 *
 *	returns map with Countries and map with Campaigns filled with Analytics Data
 */

func (parser JsonParser) Parse(fileName, aggregateType string) (*map[string]*entities.PredictionReadyEntity, error) {
	jsonFile, err := os.Open(fileName)

	if err != nil {
		log.Panicln(err)
	}

	// close file when we are done
	defer jsonFile.Close()

	jsonValue, _ := io.ReadAll(jsonFile)

	var entries []JsonAnalyticsEntry

	if err := json.Unmarshal(jsonValue, &entries); err == nil {
		result := make(map[string]*entities.PredictionReadyEntity)
		var predEntity *entities.PredictionReadyEntity

		for _, jsonEntry := range entries {
			id := jsonEntry.Country
			if aggregateType == entities.AGGREGATE_TYPE_CAMPAIGN {
				id = jsonEntry.Country
			}

			if predEntity = result[id]; predEntity == nil {
				predEntity = entities.InitPredictionReadyEntity(id)
			}

			predEntity.UsersCount += jsonEntry.Users

			fUsers := float64(jsonEntry.Users)
			//Since we are working with full numbers and not average per person - multiply average by user count
			predEntity.Ltv[0][1] += jsonEntry.Ltv1 * fUsers
			predEntity.Ltv[1][1] += jsonEntry.Ltv2 * fUsers
			predEntity.Ltv[2][1] += jsonEntry.Ltv3 * fUsers
			predEntity.Ltv[3][1] += jsonEntry.Ltv4 * fUsers
			predEntity.Ltv[4][1] += jsonEntry.Ltv5 * fUsers
			predEntity.Ltv[5][1] += jsonEntry.Ltv6 * fUsers
			predEntity.Ltv[6][1] += jsonEntry.Ltv7 * fUsers

			result[id] = predEntity
		}

		return &result, nil
	} else {
		return nil, errors.JsonFileFormatError{}
	}
}

/**
 * Storage struct to unmarshall JSON
 */
type JsonAnalyticsEntry struct {
	CampaignId string  `json:"CampaignId"`
	Country    string  `json:"Country"`
	Ltv1       float64 `json:"Ltv1"`
	Ltv2       float64 `json:"Ltv2"`
	Ltv3       float64 `json:"Ltv3"`
	Ltv4       float64 `json:"Ltv4"`
	Ltv5       float64 `json:"Ltv5"`
	Ltv6       float64 `json:"Ltv6"`
	Ltv7       float64 `json:"Ltv7"`
	Users      int     `json:"Users"`
}
