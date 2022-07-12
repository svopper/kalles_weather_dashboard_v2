package models

type DMIObservation struct {
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

type Link struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Title string `json:"title"`
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
	Name        string  `json:"name"`
}
