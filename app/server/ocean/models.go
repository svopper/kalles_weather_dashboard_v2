package ocean

type oceanViewModel struct {
	Date         string
	Observations []observation
}

type observation struct {
	stationId   string
	stationName string
	temperature float32
}
