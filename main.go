package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type WeatherObservation struct {
	Type           string    `json:"type"`
	Features       []Feature `json:"features"`
	TimeStamp      string    `json:"timeStamp"`
	NumberReturned int64     `json:"numberReturned"`
	Links          []Link    `json:"links"`
}

type Feature struct {
	Geometry   Geometry   `json:"geometry"`
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Properties struct {
	Created     string  `json:"created"`
	Observed    string  `json:"observed"`
	ParameterID string  `json:"parameterId"`
	StationID   string  `json:"stationId"`
	Value       float64 `json:"value"`
}

type Link struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type WeatherResponseData struct {
	Year int     `json:"year"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

type IndexViewModel struct {
	Date        string
	WeatherData []WeatherResponseData
}

func UnmarshalWelcome(data []byte) (WeatherObservation, error) {
	var r WeatherObservation
	err := json.Unmarshal(data, &r)
	return r, err
}

func getApiKeyFromEnvVariable() string {
	envVar := os.Getenv("DMI_API_KEY")
	if envVar == "" {
		panic("DMI_API_KEY is not set")
	}
	return envVar
}

func formatDate(date time.Time) string {
	return date.Format("2006-01-02T15:04:05Z")
}

func getWatherObservations(from, to time.Time) WeatherObservation {
	uri := fmt.Sprintf("https://dmigw.govcloud.dk/v2/metObs/collections/observation/items?datetime=%s/%s&stationId=06186&parameterId=temp_dry&bbox-crs=https://www.opengis.net/def/crs/OGC/1.3/CRS84&api-key=%s", formatDate(from), formatDate(to), getApiKeyFromEnvVariable())

	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("accept", "application/geo+json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	weatherObs, err := UnmarshalWelcome(body)

	if err != nil {
		panic(err)
	}
	return weatherObs
}

func getMinAndMax(features []Feature) (float64, float64) {
	var min, max float64
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

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		viewModel := IndexViewModel{
			Date:        time.Now().Format("January 02"),
			WeatherData: []WeatherResponseData{},
		}
		for i := 1; i <= 10; i++ {
			year := time.Now().Year() - i
			month := time.Now().Month()
			day := time.Now().Day()
			w := getWatherObservations(time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location()), time.Date(year, month, day, 23, 59, 0, 0, time.Now().Location()))
			min, max := getMinAndMax(w.Features)
			obs := WeatherResponseData{Year: year, Min: min, Max: max}
			viewModel.WeatherData = append(viewModel.WeatherData, obs)
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"data": viewModel,
		})
	})
	router.Run(":8080")
}
