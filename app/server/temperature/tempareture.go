package temperature

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/svopper/kalles_weather_dashboard_v2/app/server/util"
	"github.com/svopper/kalles_weather_dashboard_v2/app/server/util/models"
)

func getAverageMaxTemp(observations []temperatureObservation) float64 {
	var sum float64
	iterations := 0
	for _, observation := range observations {
		if observation.Max == math.Inf(1) {
			continue
		}
		sum += observation.Max
		iterations++
	}
	return util.RoundToTwoDecimal(sum / float64(iterations))
}

func getAverageMinTemp(observations []temperatureObservation) float64 {
	var sum float64
	iterations := 0
	for _, observation := range observations {
		if observation.Min == math.Inf(-1) {
			continue
		}
		sum += observation.Min
		iterations++

	}
	return util.RoundToTwoDecimal(sum / float64(iterations))
}

func getWatherObservations(from, to time.Time) models.DMIObservation {
	uri := generateTemperatureUri(from, to)
	request := util.BuildRequest(uri)
	response := util.DoRequest(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	weatherObs, err := util.UnmarshalDMIObservation(body)

	if err != nil {
		panic(err)
	}
	return weatherObs
}

func generateTemperatureUri(fromDate, toDate time.Time) string {
	uri := fmt.Sprintf(
		"https://dmigw.govcloud.dk/v2/metObs/collections/observation/items?datetime=%s/%s&stationId=06180&parameterId=temp_dry&bbox-crs=https://www.opengis.net/def/crs/OGC/1.3/CRS84&api-key=%s",
		util.FormatDate(fromDate),
		util.FormatDate(toDate),
		util.GetEnvVariable("DMI_API_KEY"),
	)
	return uri
}

func getMinAndMax(features []models.Feature) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	for _, feature := range features {
		if feature.Properties.Value < min {
			min = feature.Properties.Value
		}
		if feature.Properties.Value > max {
			max = feature.Properties.Value
		}
	}
	return min, max
}

func GetIndex(c *gin.Context) {
	viewModel := indexViewModel{
		Date:                    time.Now().Format("January 02"),
		TemperatureObservations: []temperatureObservation{},
	}
	for i := 1; i <= 10; i++ {
		year := time.Now().Year() - i
		month := time.Now().Month()
		day := time.Now().Day()
		if !util.IsLeapYear(year) && month == time.February && day == 29 {
			viewModel.TemperatureObservations = append(viewModel.TemperatureObservations, temperatureObservation{Year: year, Min: math.Inf(-1), Max: math.Inf(1)})
			continue
		}
		fromDate := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
		toDate := time.Date(year, month, day, 23, 59, 0, 0, time.Now().Location())
		w := getWatherObservations(fromDate, toDate)
		min, max := getMinAndMax(w.Features)
		obs := temperatureObservation{Year: year, Min: min, Max: max}
		viewModel.TemperatureObservations = append(viewModel.TemperatureObservations, obs)
	}
	viewModel.MaxAverage = getAverageMaxTemp(viewModel.TemperatureObservations)
	viewModel.MinAverage = getAverageMinTemp(viewModel.TemperatureObservations)
	viewModel.IsNA = func(f float64) bool { return math.IsInf(f, 0) }

	c.HTML(http.StatusOK, "index.go.tmpl", gin.H{
		"data": viewModel,
	})
}
