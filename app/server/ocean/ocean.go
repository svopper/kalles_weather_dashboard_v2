package ocean

import (
	"fmt"
	"io/ioutil"
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

func GetOcean(c *gin.Context) {
	viewModel := oceanViewModel{
		Date: time.Now().Add(-24 * time.Hour).Format("January 02"),
	}

	stationIDs := []int{util.COPENHAGEN_ID, util.HORNBAEK_ID}

	for _, stationId := range stationIDs {
		obs := getOceanObservations(stationId)
		fmt.Println(obs)

	}

	c.HTML(http.StatusOK, "ocean.go.tmpl", gin.H{
		"data": viewModel,
	})
}
