package entities

const (
	AGGREGATE_TYPE_CAMPAIGN = "campaign"
	AGGREGATE_TYPE_COUNTRY  = "country"

	PREDICTION_MODEL_LINEAR = "linear"
	PREDICTION_MODEL_QUAD   = "quad"
)

/**
 * Struct holding analytics data, with links to it's campaign and country
 */
type AnalyticsEntity struct {
	Country       *CountryEntity
	Campaign      *CampaignEntity
	AnalyticsData []float64
	UsersCount    int
}

/**
 * Country entity with name and an array of it analytics data
 */
type CountryEntity struct {
	Name      string
	Analytics []*AnalyticsEntity
}

/**
 * Campaign entity with guid and an array of it analytics data
 */
type CampaignEntity struct {
	Guid      string
	Analytics []*AnalyticsEntity
}

/**
 * Converts the campaign to a PredictionReadyEntity, summing it's analytics and users count using @fillInPredictionEntity
 */
func (campaign CampaignEntity) ConvertToPredictionReady() PredictionReadyEntity {
	result := PredictionReadyEntity{Name: campaign.Guid, UsersCount: 0, Ltv: make([][2]float64, 7)}

	fillInPredictionEntity(campaign.Analytics, &result)

	return result
}

/**
 * Converts the country to a PredictionReadyEntity, summing it's analytics and users count using @fillInPredictionEntity
 */
func (country CountryEntity) ConvertToPredictionReady() PredictionReadyEntity {
	result := PredictionReadyEntity{Name: country.Name, UsersCount: 0, Ltv: make([][2]float64, 7)}

	fillInPredictionEntity(country.Analytics, &result)

	return result
}

/**
 * Common method to convert country or campaign analytics to the Prediction Ready entity
 */
func fillInPredictionEntity(Analytics []*AnalyticsEntity, predictionReady *PredictionReadyEntity) {
	for _, analyticItem := range Analytics {
		predictionReady.UsersCount += analyticItem.UsersCount
		for i := 0; i < 7; i++ {
			predictionReady.Ltv[i] = [2]float64{float64(i + 1), predictionReady.Ltv[i][1] + analyticItem.AnalyticsData[i]}
		}
	}
}

/**
 * Convert campaigns map to an array of prediction ready values Predictor expects
 */
func NormalizeAndSumValuesCampaign(items *map[string]*CampaignEntity) []PredictionReadyEntity {
	result := make([]PredictionReadyEntity, 0)
	for _, campaign := range *items {

		result = append(result, campaign.ConvertToPredictionReady())
	}

	return result
}

/**
 * Convert countries map to an array of prediction ready values Predictor expects
 */
func NormalizeAndSumValuesCountry(items *map[string]*CountryEntity) []PredictionReadyEntity {
	result := make([]PredictionReadyEntity, 0)
	for _, country := range *items {
		result = append(result, country.ConvertToPredictionReady())
	}

	return result
}

func InitCountryEntity(name string) *CountryEntity {
	m := make([]*AnalyticsEntity, 0)
	country := CountryEntity{Name: name, Analytics: m}
	return &country
}

func InitCampaignEntity(guid string) *CampaignEntity {
	m := make([]*AnalyticsEntity, 0)
	campaign := CampaignEntity{Guid: guid, Analytics: m}
	return &campaign
}

func InitAnalyticsEntity(country *CountryEntity, campaign *CampaignEntity) *AnalyticsEntity {
	analtyics := AnalyticsEntity{country, campaign, make([]float64, 7), 0}
	country.Analytics = append(country.Analytics, &analtyics)
	campaign.Analytics = append(campaign.Analytics, &analtyics)

	return &analtyics
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
