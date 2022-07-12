package station_map

type mapViewModel struct {
	Stations []stationViewModel
}

type stationViewModel struct {
	Coordinates coordinate
	Name        string
	StationID   string
}

type coordinate struct {
	Lat float64
	Lng float64
}
