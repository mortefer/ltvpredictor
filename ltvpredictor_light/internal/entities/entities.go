package entities

const (
	AGGREGATE_TYPE_CAMPAIGN = "campaign"
	AGGREGATE_TYPE_COUNTRY  = "country"

	PREDICTION_MODEL_LINEAR = "linear"
	PREDICTION_MODEL_QUAD   = "quad"
)

func InitPredictionReadyEntity(id string) *PredictionReadyEntity {
	res := PredictionReadyEntity{Name: id, Ltv: make([][2]float64, 7), UsersCount: 0}
	res.Ltv[0][0] = 1
	res.Ltv[1][0] = 2
	res.Ltv[2][0] = 3
	res.Ltv[3][0] = 4
	res.Ltv[4][0] = 5
	res.Ltv[5][0] = 6
	res.Ltv[6][0] = 7

	return &res
}

/**
 * Helper struct holding information formatted for the prediction library
 */
type PredictionReadyEntity struct {
	Name       string
	Ltv        [][2]float64
	UsersCount int
}

/**
 * An object to be passed back when prediction is done, containing name to print and the prediction value (accuracy value is for the debugging)
 */
type PredictionResult struct {
	Name     string
	Value    float64
	Accuracy float64
}
