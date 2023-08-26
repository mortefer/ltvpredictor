package predictors

import (
	"log"
	"mortefer/ltvpredictor/internal/entities"
	"os"
	"strconv"

	"github.com/Konstantin8105/regression"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func StartPrediction(items []entities.PredictionReadyEntity, predictionModel string, outputImage int, doneChannel chan []entities.PredictionResult) {

	predictionChan := make(chan entities.PredictionResult)

	/**
	 *	Launch predictions routines
	 */
	for _, item := range items {
		go performPrediction(item, predictionModel, outputImage, predictionChan)
	}

	predictionResults := make([]entities.PredictionResult, 0)

	/**
	 * Read our predictions from the channel, form the output array and write to main channel
	 */
	for i := 0; i < len(items); i++ {
		predictionResults = append(predictionResults, <-predictionChan)
	}

	doneChannel <- predictionResults
}

func performPrediction(predictionReady entities.PredictionReadyEntity, predictionModel string, outputVisual int, resultChan chan entities.PredictionResult) {
	var predictor IPredictor

	/**
	 * Init predictior here since it is not thread safe to init it in the calling method above and the pass here a reference,
	 * and I don't like mutexes :)
	 */
	switch predictionModel {
	case entities.PREDICTION_MODEL_LINEAR:
		predictor = &LinearPredictor{}
	case entities.PREDICTION_MODEL_QUAD:
		predictor = &QuadraticPredictor{}
	default:
		predictor = &QuadraticPredictor{}

	}

	var predictionResult float64

	err := predictor.Train(predictionReady.Ltv)

	if err == nil {
		if outputVisual != 0 {
			/**
			* Since we are running in a separate routine, let's have fun and predict all the days up to day 60
			* with that we are able to build the image
			 */
			for i := 1.0; i <= 53; i++ {
				predictionResult = predictor.Predict(i + 7) // +7 because we already have 7 days of data
				predictionReady.Ltv = append(predictionReady.Ltv, [2]float64{i + 7, predictionResult})
			}

			generateVisual(&predictionReady, predictionModel)
		} else {
			predictionResult = predictor.Predict(60)
		}

		//we need to store just the data for one user, so devide by user amount
		predictionResult = predictionResult / float64(predictionReady.UsersCount)

		resultChan <- entities.PredictionResult{Value: predictionResult, Name: predictionReady.Name, Accuracy: predictor.GetAccuracy()}
	} else {
		resultChan <- entities.PredictionResult{Value: -1, Name: predictionReady.Name + " - ERROR"}
	}
}

/**
 * Small bonus method to generate charts to view the generated data
 */
func generateVisual(data *entities.PredictionReadyEntity, nameAppend string) {
	lineChart := charts.NewLine()

	lineChart.SetGlobalOptions(
		//charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeVintage}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Income for " + data.Name + " per user",
			Subtitle: "Total user count: " + strconv.Itoa(data.UsersCount),
		}))

	realPoints := make([]opts.LineData, 0)
	predictedPoints := make([]opts.LineData, 0)
	titles := make([]string, 0)
	for num, item := range data.Ltv {
		if num < 7 {
			realPoints = append(realPoints, opts.LineData{Value: item[1] / float64(data.UsersCount)})
			predictedPoints = append(predictedPoints, opts.LineData{Value: item[1] / float64(data.UsersCount)})
		} else {
			predictedPoints = append(predictedPoints, opts.LineData{Value: item[1] / float64(data.UsersCount)})
		}
		titles = append(titles, strconv.Itoa(int(item[0])))
	}

	lineChart.SetXAxis(titles).
		AddSeries("Predicted data", predictedPoints).
		AddSeries("Actual data", realPoints)

	f, err := os.Create("charts/" + data.Name + "-" + nameAppend + ".html")
	if err == nil {
		lineChart.Render(f)
	} else {
		log.Fatalln(err)
	}
}

/**
 * Iterface for the linear and cubic preductors wrappers
 * Train - set the data and get the koeffs of the function
 * Predict - calculate Y value of the function
 */
type IPredictor interface {
	Train([][2]float64) error
	Predict(float64) float64
	GetAccuracy() float64
}

/**
 * Linear extrapolator wrapper
 */
type LinearPredictor struct {
	a, b, R2 float64
}

/**
 * Train the model with the existing data
 */
func (linear *LinearPredictor) Train(inputMatrix [][2]float64) error {
	a, b, R2, err := regression.Linear(inputMatrix)
	if err == nil {
		linear.a = a
		linear.b = b
		linear.R2 = R2
		return nil
	} else {
		return err
	}
}

/**
 * Run the linear prediction function
 */
func (linear LinearPredictor) Predict(x float64) float64 {
	return linear.a*x + linear.b
}

/**
 * Return accuracy
 */
func (linear LinearPredictor) GetAccuracy() float64 {
	return linear.R2
}

/**
 * Qadratic extrapolator wrapper
 */
type QuadraticPredictor struct {
	a, b, c, R2 float64
}

/**
 * Train the quadratic model with the existing data
 */
func (quad *QuadraticPredictor) Train(inputMatrix [][2]float64) error {
	a, b, c, R2, err := regression.Quadratic(inputMatrix)
	if err == nil {
		quad.a = a
		quad.b = b
		quad.c = c
		quad.R2 = R2
		return nil
	} else {
		return err
	}
}

/**
 * Run the quadratic prediction function
 */
func (quad QuadraticPredictor) Predict(x float64) float64 {
	return quad.a*x*x + quad.b*x + quad.c
}

/**
 * Return accuracy
 */
func (quad QuadraticPredictor) GetAccuracy() float64 {
	return quad.R2
}
