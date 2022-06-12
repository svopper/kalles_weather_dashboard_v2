package temperature

type indexViewModel struct {
	Date                    string
	TemperatureObservations []temperatureObservation
	MaxAverage              float64
	MinAverage              float64
	IsNA                    func(float64) bool
}

type temperatureObservation struct {
	Year int     `json:"year"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
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

type link struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type properties struct {
	Created     string  `json:"created"`
	Observed    string  `json:"observed"`
	ParameterID string  `json:"parameterId"`
	StationID   string  `json:"stationId"`
	Value       float64 `json:"value"`
}
