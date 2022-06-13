package ocean

type oceanViewModel struct {
	Date         string
	Observations []observation
	IsNA         func(float64) bool
}

type observation struct {
	StationId   int
	StationName string
	Temperature float64
}
