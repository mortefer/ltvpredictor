package main

import (
	"flag"
	"log"
	"mortefer/ltvpredictor_light/internal/entities"
	"mortefer/ltvpredictor_light/internal/parsers"
	"mortefer/ltvpredictor_light/internal/predictors"
)

func main() {
	/**
	 * Init the flags that we are expecting, set the default values so that utility runs out of the box
	 */
	predictionModel := flag.String("model", entities.PREDICTION_MODEL_LINEAR, "Prediction model switch, valid values are '"+entities.PREDICTION_MODEL_QUAD+"', '"+entities.PREDICTION_MODEL_LINEAR+"'")
	sourceFile := flag.String("source", "test_data.json", "Source file name, existing in the data folder")
	aggregateType := flag.String("aggregate", "country", "Aggregate type, either '"+entities.AGGREGATE_TYPE_COUNTRY+"' or '"+entities.AGGREGATE_TYPE_CAMPAIGN+"'")

	flag.Parse()

	fileName, extension, err := parsers.CheckInputIsValidAndFileData(*aggregateType, *predictionModel, *sourceFile)
	if err != nil {
		//TODO: we can implement different error handling basing on returned error type
		log.Fatal(err)
	} else {
		//init parser by file type
		//TODO: I guess I could have used the reflection package here, but this would have taken me quite some time to explore :)
		parser := parsers.InitParserByType(extension)

		/**
		 * Here goes the simplification.
		 * We don't pay attention to many to many relationship, we just take the data that we need, either grouping by campaign or country
		 * Having taken the identifier, we just sum up the data for it
		 */
		predictionReadyValues, err := parser.Parse(fileName, *aggregateType)
		if err == nil {
			doneChan := make(chan []entities.PredictionResult)

			go predictors.StartPrediction(predictionReadyValues, *predictionModel, doneChan)

			//wait until subroutines are done and read the results
			results := <-doneChan

			for _, predRes := range results {
				//we can check for -1 value here to see if error happened anywhere
				//Log received data according to task
				log.Println(predRes.Name, predRes.Value)
			}

			log.Println("Done")
		} else {
			log.Fatal(err)
		}
	}
}
