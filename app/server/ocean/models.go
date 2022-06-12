package ocean

type oceanViewModel struct {
	Date         string
	Observations []observation
}

type observation struct {
	stationId   int
	stationName string
	temperature float64
}
