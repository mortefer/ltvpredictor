package main

import (
	"flag"
	"log"
	"mortefer/ltvpredictor/internal/entities"
	"mortefer/ltvpredictor/internal/parsers"
	"mortefer/ltvpredictor/internal/predictors"
)

func main() {
	/**
	 * Init the flags that we are expecting, set the default values so that utility runs out of the box
	 */
	predictionModel := flag.String("model", entities.PREDICTION_MODEL_QUAD, "Prediction model switch, valid values are '"+entities.PREDICTION_MODEL_QUAD+"', '"+entities.PREDICTION_MODEL_LINEAR+"'")
	sourceFile := flag.String("source", "test_data.json", "Source file name, existing in the data folder")
	aggregateType := flag.String("aggregate", "country", "Aggregate type, either '"+entities.AGGREGATE_TYPE_COUNTRY+"' or '"+entities.AGGREGATE_TYPE_CAMPAIGN+"'")
	outputImage := flag.Int("graph", 1, "Should the software try and generate graph charts basing on prediction data, valid values are: 0 or 1")

	flag.Parse()

	fileName, extension, err := parsers.CheckInputIsValidAndFileData(*aggregateType, *predictionModel, *sourceFile)
	if err != nil {
		//TODO: we can implement different error handling basing on returned error type
		log.Fatal(err)
	} else {
		/**
		 * init parser by file type
		 * TODO: I guess I could have used the reflection package here, but this would have taken me quite some time to explore :)
		 */
		parser := parsers.InitParserByType(extension)

		/**
		 * We return both countries and campaigns here, which might seem like an overkill, but gives us the ability to extend
		 * the functionality, we can play wiht both data sets, build graphs whatever we like, or even make cross referenses
		 */
		countries, campaigns, err := parser.Parse(fileName)
		if err == nil {
			/**
			 * same here, Countries and Campaigns are not that much different, and reflection could have been used to parse them
			 * OR I could have used a common interface for Campaign and Country, but both of those are an overkill for just one method.
			 * Simplier solution is in the Simplified project
			 */
			var predictionReadyValues []entities.PredictionReadyEntity

			if *aggregateType == entities.AGGREGATE_TYPE_CAMPAIGN {
				predictionReadyValues = entities.NormalizeAndSumValuesCampaign(campaigns)
			} else {
				predictionReadyValues = entities.NormalizeAndSumValuesCountry(countries)
			}

			doneChan := make(chan []entities.PredictionResult)

			go predictors.StartPrediction(predictionReadyValues, *predictionModel, *outputImage, doneChan)

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
