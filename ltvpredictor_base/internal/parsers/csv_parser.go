package parsers

import (
	"encoding/csv"
	"log"
	"mortefer/ltvpredictor/internal/entities"
	"mortefer/ltvpredictor/internal/errors"
	"os"
	"strconv"
)

/**
 *	CsvParser, implements InputFileParser interface
 */
type CsvParser struct {

	/** map to store Country information, key is country name */
	Countries map[string]*entities.CountryEntity

	/** map to store Campaigns information, key is campaign GUID */
	Campaigns map[string]*entities.CampaignEntity

	/** map to hold analytics data, key is CountryName-CampaignGUID */
	Analytics map[string]*entities.AnalyticsEntity
}

/**
 *	Parse supplied CSV file, while trying to check for validity of CSV data
 *  Performs the summing of the user data values into analytics data, keeping in mind the day user became customer
 *  Receives fileName string
 *
 *	returns map with Countries and map with Campaigns filled with Analytics Data
 */

func (parser CsvParser) Parse(fileName string) (*map[string]*entities.CountryEntity, *map[string]*entities.CampaignEntity, error) {
	csvFile, err := os.Open(fileName)

	if err != nil {
		log.Panicln(err)
	}
	// close file when we are done
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	//skip header, while checking that the format is ok
	csvEntry, err := csvReader.Read()

	if err != nil {
		return nil, nil, errors.CsvFileFormatError{}
	}

	/**
	 * check that the format is the one expected
	 * for the sake of test task we just check for the expected length,
	 * otherwise we could have checked the header to be valid
	 */
	if len(csvEntry) != 10 {
		return nil, nil, errors.CsvFileFormatError{}
	}

	var analytics *entities.AnalyticsEntity
	var country *entities.CountryEntity
	var campaign *entities.CampaignEntity

	//tmp vars, moved out of the loop to preserve memory
	var offset int
	var value float64

	for {
		csvEntry, err = csvReader.Read()
		if err != nil {
			break
		}
		if campaign = parser.Campaigns[csvEntry[1]]; campaign == nil {
			campaign = entities.InitCampaignEntity(csvEntry[1])
			parser.Campaigns[csvEntry[1]] = campaign
		}

		if country = parser.Countries[csvEntry[2]]; country == nil {
			country = entities.InitCountryEntity(csvEntry[2])
			parser.Countries[csvEntry[2]] = country
		}

		/**
		 * Here we will definitely have several occurences of Country-Campaigns records, so we need to sum the values
		 */
		if analytics = parser.Analytics[country.Name+"-"+campaign.Guid]; analytics == nil {
			analytics = entities.InitAnalyticsEntity(country, campaign)
			parser.Analytics[country.Name+"-"+campaign.Guid] = analytics
		}

		/**
		 * Set the offset so we sum correct LTV Values
		 */
		offset = 0
		for i := 9; i >= 6; i-- {
			if csvEntry[i] == "0" {
				offset++
			} else {
				break
			}
		}

		for i := offset; i < 7; i++ {
			value, _ = strconv.ParseFloat(csvEntry[3+i-offset], 32)
			analytics.AnalyticsData[i] += value
		}
		analytics.UsersCount++

	}

	return &parser.Countries, &parser.Campaigns, nil
}
