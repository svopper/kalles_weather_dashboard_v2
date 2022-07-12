package station_map

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/svopper/kalles_weather_dashboard_v2/app/server/util"
	"github.com/svopper/kalles_weather_dashboard_v2/app/server/util/models/station"
)

func generateStationsUri() string {
	uri := fmt.Sprintf(
		"https://dmigw.govcloud.dk/v2/oceanObs/collections/station/items?status=Active&bbox-crs=https://www.opengis.net/def/crs/OGC/1.3/CRS84&api-key=%s",
		util.GetEnvVariable("DMI_OCEAN_API_KEY"),
	)
	return uri
}

func getStationsObservations() station.Stations {
	uri := generateStationsUri()
	request := util.BuildRequest(uri)
	response := util.DoRequest(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	stations, err := station.UnmarshalStation(body)

	if err != nil {
		panic(err)
	}
	return stations
}

func GetMap(c *gin.Context) {
	stations := getStationsObservations()
	var stationsParsed []stationViewModel

	for _, s := range stations.Features {
		sParsed := stationViewModel{
			Coordinates: coordinate{
				Lat: s.Geometry.Coordinates[1],
				Lng: s.Geometry.Coordinates[0],
			},
			Name:      s.Properties.Name,
			StationID: s.Properties.StationID,
		}
		stationsParsed = append(stationsParsed, sParsed)
	}

	c.HTML(http.StatusOK, "map.go.tmpl", gin.H{
		"ocean_stations":   stationsParsed,
		"weather_stations": nil,
	})
}
