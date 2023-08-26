package parsers

import (
	"encoding/json"
	"io"
	"log"
	"mortefer/ltvpredictor/internal/entities"
	"mortefer/ltvpredictor/internal/errors"
	"os"
)

/**
 *	JsonParser, implements InputFileParser interface
 */
type JsonParser struct {

	/** map to store Country information, key is country name */
	Countries map[string]*entities.CountryEntity

	/** map to store Campaigns information, key is campaign GUID */
	Campaigns map[string]*entities.CampaignEntity

	/** map to hold analytics data, key is CountryName-CampaignGUID */
	Analytics map[string]*entities.AnalyticsEntity
}

/**
 *	Parse supplied JSON file, while checking for validity of JSON data
 *  Receives fileName string
 *
 *	returns map with Countries and map with Campaigns filled with Analytics Data
 */

func (parser JsonParser) Parse(fileName string) (*map[string]*entities.CountryEntity, *map[string]*entities.CampaignEntity, error) {
	jsonFile, err := os.Open(fileName)

	if err != nil {
		log.Panicln(err)
	}

	// close file when we are done
	defer jsonFile.Close()

	jsonValue, _ := io.ReadAll(jsonFile)

	var entries []JsonAnalyticsEntry

	if err := json.Unmarshal(jsonValue, &entries); err == nil {
		var analytics *entities.AnalyticsEntity
		var country *entities.CountryEntity
		var campaign *entities.CampaignEntity

		for _, jsonEntry := range entries {
			if country = parser.Countries[jsonEntry.Country]; country == nil {
				country = entities.InitCountryEntity(jsonEntry.Country)
				parser.Countries[jsonEntry.Country] = country
			}

			if campaign = parser.Campaigns[jsonEntry.CampaignId]; campaign == nil {
				campaign = entities.InitCampaignEntity(jsonEntry.CampaignId)
				parser.Campaigns[jsonEntry.CampaignId] = campaign
			}

			/**
			 * This is mostly a failsafe mechanism in case we have more then one unique entry for Country-Campaign occurence
			 * This is not the case for the test JSON data, which seems straight forward, however I though this is a good idea to have
			 */
			if analytics = parser.Analytics[jsonEntry.Country+"-"+jsonEntry.CampaignId]; analytics == nil {
				analytics = entities.InitAnalyticsEntity(country, campaign)
				parser.Analytics[jsonEntry.Country+"-"+jsonEntry.CampaignId] = analytics
			}
			analytics.UsersCount += jsonEntry.Users

			fUsers := float64(jsonEntry.Users)
			//Since we are working with full numbers and not average per person - multiply average by user count
			analytics.AnalyticsData[0] += jsonEntry.Ltv1 * fUsers
			analytics.AnalyticsData[1] += jsonEntry.Ltv2 * fUsers
			analytics.AnalyticsData[2] += jsonEntry.Ltv3 * fUsers
			analytics.AnalyticsData[3] += jsonEntry.Ltv4 * fUsers
			analytics.AnalyticsData[4] += jsonEntry.Ltv5 * fUsers
			analytics.AnalyticsData[5] += jsonEntry.Ltv6 * fUsers
			analytics.AnalyticsData[6] += jsonEntry.Ltv7 * fUsers
		}

		return &parser.Countries, &parser.Campaigns, nil
	} else {
		return nil, nil, errors.JsonFileFormatError{}
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
