package ocean

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

func generateOceanUri(stationId int) string {
	uri := fmt.Sprintf(
		"https://dmigw.govcloud.dk/v2/oceanObs/collections/observation/items?period=latest-day&stationId=%d&parameterId=tw&sortorder=observed,DESC&bbox-crs=https://www.opengis.net/def/crs/OGC/1.3/CRS84&api-key=%s",
		stationId,
		util.GetEnvVariable("DMI_OCEAN_API_KEY"),
	)
	return uri
}

func getOceanObservations(stationId int) models.DMIObservation {
	uri := generateOceanUri(stationId)
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

func getMax(features []models.Feature) float64 {
	max := math.Inf(-1)
	for _, feature := range features {
		if feature.Properties.Value > max {
			max = feature.Properties.Value
		}
	}
	return max
}

func GetOcean(c *gin.Context) {
	viewModel := oceanViewModel{
		Date: time.Now().Format("January 02"), // get observation for yesterday
		IsNA: func(f float64) bool { return math.IsInf(f, 0) },
	}

	for stationId, stationName := range util.OCEAN_STATION_MAP {
		obs := getOceanObservations(stationId)
		observation := observation{
			StationId:   stationId,
			StationName: stationName,
			Temperature: getMax(obs.Features),
		}
		viewModel.Observations = append(viewModel.Observations, observation)
	}

	c.HTML(http.StatusOK, "ocean.go.tmpl", gin.H{
		"data": viewModel,
	})
}
