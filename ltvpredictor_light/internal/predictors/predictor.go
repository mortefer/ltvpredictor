package predictors

import (
	"mortefer/ltvpredictor_light/internal/entities"

	"github.com/Konstantin8105/regression"
)

func StartPrediction(items *map[string]*entities.PredictionReadyEntity, predictionModel string, doneChannel chan []entities.PredictionResult) {

	predictionChan := make(chan entities.PredictionResult)

	/**
	 *	Launch predictions routines
	 */
	for _, item := range *items {
		go performPrediction(item, predictionModel, predictionChan)
	}

	predictionResults := make([]entities.PredictionResult, 0)

	/**
	 * Read our predictions from the channel, form the output array and write to main channel
	 */
	for range *items {
		predictionResults = append(predictionResults, <-predictionChan)
	}

	doneChannel <- predictionResults
}

func performPrediction(predictionReady *entities.PredictionReadyEntity, predictionModel string, resultChan chan entities.PredictionResult) {
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

		predictionResult = predictor.Predict(60)
		//we need to store just the data for one user, so devide by user amount
		predictionResult = predictionResult / float64(predictionReady.UsersCount)

		resultChan <- entities.PredictionResult{Value: predictionResult, Name: predictionReady.Name, Accuracy: predictor.GetAccuracy()}
	} else {
		resultChan <- entities.PredictionResult{Value: -1, Name: predictionReady.Name + " - ERROR"}
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
