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
