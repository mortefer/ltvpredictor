package parsers

import (
	"encoding/csv"
	"log"
	"mortefer/ltvpredictor_light/internal/entities"
	"mortefer/ltvpredictor_light/internal/errors"
	"os"
	"strconv"
)

/**
 *	CsvParser, implements InputFileParser interface
 */
type CsvParser struct{}

/**
 *	Parse supplied CSV file, while trying to check for validity of CSV data
 *  Performs the summing of the user data values into analytics data, keeping in mind the day user became customer
 *  Receives fileName string
 *
 *	returns map with Countries and map with Campaigns filled with Analytics Data
 */

func (parser CsvParser) Parse(fileName, aggregateType string) (*map[string]*entities.PredictionReadyEntity, error) {
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
		return nil, errors.CsvFileFormatError{}
	}

	/**
	 * check that the format is the one expected
	 * for the sake of test task we just check for the expected length,
	 * otherwise we could have checked the header to be valid
	 */
	if len(csvEntry) != 10 {
		return nil, errors.CsvFileFormatError{}
	}

	//tmp vars, moved out of the loop to preserve memory
	var offset int
	var value float64

	result := make(map[string]*entities.PredictionReadyEntity)
	var predEntity *entities.PredictionReadyEntity

	for {
		csvEntry, err = csvReader.Read()
		if err != nil {
			break
		}

		id := csvEntry[1] //campaign name
		if aggregateType == entities.AGGREGATE_TYPE_COUNTRY {
			id = csvEntry[2]
		}

		if predEntity = result[id]; predEntity == nil {
			predEntity = entities.InitPredictionReadyEntity(id)
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
			predEntity.Ltv[i][1] += value
		}

		predEntity.UsersCount++

		result[id] = predEntity

	}

	return &result, nil
}
