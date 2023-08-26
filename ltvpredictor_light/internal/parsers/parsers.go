package parsers

import (
	"os"
	"path/filepath"

	"mortefer/ltvpredictor_light/internal/entities"
	"mortefer/ltvpredictor_light/internal/errors"
)

/**
 * Interface for the Parsers implementation
 */

type IInputFileParser interface {
	Parse(string, string) (*map[string]*entities.PredictionReadyEntity, error)
}

/**
 * Parse input parameters and validate them for expected values
 *
 * return file extension information for future parser initialization
 */
func CheckInputIsValidAndFileData(aggregateType, predictionModel, inputSourceFile string) (string, string, error) {

	if aggregateType == "" || (aggregateType != entities.AGGREGATE_TYPE_CAMPAIGN && aggregateType != entities.AGGREGATE_TYPE_COUNTRY) {
		return "", "", errors.AggregateTypeIncorrectError{}
	}

	if predictionModel == "" || (predictionModel != entities.PREDICTION_MODEL_QUAD && predictionModel != entities.PREDICTION_MODEL_LINEAR) {
		return "", "", errors.PredictionModelIncorrectError{}
	}

	fileName := "data/" + inputSourceFile

	file, err := os.Stat(fileName)
	if err == nil {
		ext := filepath.Ext(fileName)

		if !file.IsDir() && (ext == ".csv" || ext == ".json") {
			return fileName, ext, nil
		} else {
			return "", "", errors.InputFileError{Message: "Specified file exists, but the format is not supported. It should be either csv or json."}
		}
	} else {
		return "", "", errors.InputFileError{Message: "Unable to open source file: 'data/" + inputSourceFile + "', check that the file exists"}
	}
}

/**
 * return Parser object basing on file type, no need to check for validity of extension, we did that before
 */
func InitParserByType(ext string) IInputFileParser {
	if ext == ".json" {
		return JsonParser{}
	}

	return CsvParser{}
}
