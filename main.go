package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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

type TemperatureObservation struct {
	Year int     `json:"year"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

type IndexViewModel struct {
	Date                    string
	TemperatureObservations []TemperatureObservation
	MaxAverage              float64
	MinAverage              float64
	IsNA                    func(float64) bool
}

func UnmarshalWeatherObservation(data []byte) (WeatherObservation, error) {
	var r WeatherObservation
	err := json.Unmarshal(data, &r)
	return r, err
}

func getEnvVariable(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is not set", name))
	}
	return value
}

func formatDate(date time.Time) string {
	return date.Format("2006-01-02T15:04:05Z")
}

func generateUri(fromDate, toDate time.Time) string {
	uri := fmt.Sprintf(
		"https://dmigw.govcloud.dk/v2/metObs/collections/observation/items?datetime=%s/%s&stationId=06186&parameterId=temp_dry&bbox-crs=https://www.opengis.net/def/crs/OGC/1.3/CRS84&api-key=%s",
		formatDate(fromDate),
		formatDate(toDate),
		getEnvVariable("DMI_API_KEY"),
	)
	return uri
}

func buildRequest(uri string) *http.Request {
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("accept", "application/geo+json")
	return req
}

func doRequest(request *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	return resp
}

func getWatherObservations(from, to time.Time) WeatherObservation {
	uri := generateUri(from, to)
	request := buildRequest(uri)
	response := doRequest(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	weatherObs, err := UnmarshalWeatherObservation(body)

	if err != nil {
		panic(err)
	}
	return weatherObs
}

func getMinAndMax(features []Feature) (float64, float64) {
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

func roundToTwoDecimal(num float64) float64 {
	rounded := math.Round(num*10) / 10
	return rounded
}

func getAverageMaxTemp(observations []TemperatureObservation) float64 {
	var sum float64
	iterations := 0
	for _, observation := range observations {
		if observation.Max == math.Inf(1) {
			continue
		}
		sum += observation.Max
		iterations++
	}
	return roundToTwoDecimal(sum / float64(iterations))
}

func getAverageMinTemp(observations []TemperatureObservation) float64 {
	var sum float64
	iterations := 0
	for _, observation := range observations {
		if observation.Min == math.Inf(-1) {
			continue
		}
		sum += observation.Min
		iterations++

	}
	return roundToTwoDecimal(sum / float64(iterations))
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func isFeb29(year int, month time.Month, day int) bool {
	return isLeapYear(year) && month == time.February && day == 29
}

func isNA(number float64) bool {
	return math.IsInf(number, 0)
}

func getIndex(c *gin.Context) {
	viewModel := IndexViewModel{
		Date:                    time.Now().Format("January 02"),
		TemperatureObservations: []TemperatureObservation{},
	}
	for i := 1; i <= 10; i++ {
		year := time.Now().Year() - i
		month := time.Now().Month()
		day := time.Now().Day()
		if !isLeapYear(year) && month == time.February && day == 29 {
			viewModel.TemperatureObservations = append(viewModel.TemperatureObservations, TemperatureObservation{Year: year, Min: math.Inf(-1), Max: math.Inf(1)})
			continue
		}
		fromDate := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
		toDate := time.Date(year, month, day, 23, 59, 0, 0, time.Now().Location())
		w := getWatherObservations(fromDate, toDate)
		min, max := getMinAndMax(w.Features)
		obs := TemperatureObservation{Year: year, Min: min, Max: max}
		viewModel.TemperatureObservations = append(viewModel.TemperatureObservations, obs)
	}
	viewModel.MaxAverage = getAverageMaxTemp(viewModel.TemperatureObservations)
	viewModel.MinAverage = getAverageMinTemp(viewModel.TemperatureObservations)
	viewModel.IsNA = isNA

	c.HTML(http.StatusOK, "index.go.tmpl", gin.H{
		"data": viewModel,
	})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/", getIndex)
	router.Run(":8080")
}
