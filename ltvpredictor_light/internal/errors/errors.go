package errors

type AggregateTypeIncorrectError struct{}

func (err AggregateTypeIncorrectError) Error() string {
	return "Specified Aggregate type is incorrect. Run with -h parameter for help."
}

type PredictionModelIncorrectError struct{}

func (err PredictionModelIncorrectError) Error() string {
	return "Specified Prediction Model is incorrect. Run with -h parameter for help."
}

type InputFileError struct {
	Message string
}

func (err InputFileError) Error() string {
	return err.Message
}

type JsonFileFormatError struct{}

func (err JsonFileFormatError) Error() string {
	return "Cannot parse specified JSON file - format is unknown."
}

type CsvFileFormatError struct{}

func (err CsvFileFormatError) Error() string {
	return "Cannot parse specified CSV file - format is unknown."
}
