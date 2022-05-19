package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type weatherObservation struct {
	Type           string    `json:"type"`
	Features       []feature `json:"features"`
	TimeStamp      string    `json:"timeStamp"`
	NumberReturned int64     `json:"numberReturned"`
	Links          []link    `json:"links"`
}

type feature struct {
	Geometry   Geometry   `json:"geometry"`
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Properties properties `json:"properties"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type properties struct {
	Created     string  `json:"created"`
	Observed    string  `json:"observed"`
	ParameterID string  `json:"parameterId"`
	StationID   string  `json:"stationId"`
	Value       float64 `json:"value"`
}

type link struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type temperatureObservation struct {
	Year int     `json:"year"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

type indexViewModel struct {
	Date                    string
	TemperatureObservations []temperatureObservation
	MaxAverage              float64
	MinAverage              float64
	IsNA                    func(float64) bool
}

func unmarshalWeatherObservation(data []byte) (weatherObservation, error) {
	var r weatherObservation
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		buf := new(strings.Builder)
		io.Copy(buf, resp.Body)
		panic(fmt.Sprintf("Request failed with status code %d. Error: %s", resp.StatusCode, buf.String()))
	}
	if err != nil {
		panic(err)
	}

	return resp
}

func getWatherObservations(from, to time.Time) weatherObservation {
	uri := generateUri(from, to)
	request := buildRequest(uri)
	response := doRequest(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	weatherObs, err := unmarshalWeatherObservation(body)

	if err != nil {
		panic(err)
	}
	return weatherObs
}

func getMinAndMax(features []feature) (float64, float64) {
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
	return roundToTwoDecimal(sum / float64(iterations))
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
	viewModel := indexViewModel{
		Date:                    time.Now().Format("January 02"),
		TemperatureObservations: []temperatureObservation{},
	}
	for i := 1; i <= 10; i++ {
		year := time.Now().Year() - i
		month := time.Now().Month()
		day := time.Now().Day()
		if !isLeapYear(year) && month == time.February && day == 29 {
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
	viewModel.IsNA = isNA

	c.HTML(http.StatusOK, "index.go.tmpl", gin.H{
		"data": viewModel,
	})
}

func InstantiateControllers() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.HTML(http.StatusInternalServerError, "error.go.tmpl", gin.H{
				"error": err,
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.LoadHTMLGlob("app/server/templates/*")
	router.Static("/assets", "app/server/assets")

	router.GET("/", getIndex)

	return router
}
