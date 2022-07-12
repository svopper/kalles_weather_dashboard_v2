// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    welcome, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()

package station

import "encoding/json"

func UnmarshalStation(data []byte) (Stations, error) {
	var r Stations
	err := json.Unmarshal(data, &r)
	return r, err
}

// func (r *Stations) Marshal() ([]byte, error) {
// 	return json.Marshal(r)
// }

type Stations struct {
	Type           string           `json:"type"`
	Features       []FeatureElement `json:"features"`
	TimeStamp      string           `json:"timeStamp"`
	NumberReturned int64            `json:"numberReturned"`
	Links          []Link           `json:"links"`
}

type FeatureElement struct {
	Geometry   Geometry    `json:"geometry"`
	ID         string      `json:"id"`
	Type       FeatureType `json:"type"`
	Properties Properties  `json:"properties"`
}

type Geometry struct {
	Coordinates []float64    `json:"coordinates"`
	Type        GeometryType `json:"type"`
}

type Properties struct {
	BarometerHeight *float64       `json:"barometerHeight"`
	Country         Country        `json:"country"`
	Created         string         `json:"created"`
	Name            string         `json:"name"`
	OperationFrom   string         `json:"operationFrom"`
	OperationTo     *string        `json:"operationTo"`
	Owner           Owner          `json:"owner"`
	ParameterID     []ParameterID  `json:"parameterId"`
	RegionID        *string        `json:"regionId"`
	StationHeight   *float64       `json:"stationHeight"`
	StationID       string         `json:"stationId"`
	Status          Status         `json:"status"`
	Type            PropertiesType `json:"type"`
	Updated         interface{}    `json:"updated"`
	ValidFrom       string         `json:"validFrom"`
	ValidTo         *string        `json:"validTo"`
	WmoCountryCode  *string        `json:"wmoCountryCode"`
	WmoStationID    *string        `json:"wmoStationId"`
}

type Link struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type GeometryType string

const (
	Point GeometryType = "Point"
)

type Country string

const (
	Denmark   Country = "DNK"
	Greenland Country = "GRL"
)

type Owner string

const (
	DMI             Owner = "DMI"
	DanskeLufthavne Owner = "Danske lufthavne"
	HavneKommunerMv Owner = "Havne, Kommuner mv."
)

type ParameterID string

const (
	CloudCover            ParameterID = "cloud_cover"
	CloudHeight           ParameterID = "cloud_height"
	Humidity              ParameterID = "humidity"
	HumidityPast1H        ParameterID = "humidity_past1h"
	LeavHumDurPast10Min   ParameterID = "leav_hum_dur_past10min"
	LeavHumDurPast1H      ParameterID = "leav_hum_dur_past1h"
	PrecipDurPast10Min    ParameterID = "precip_dur_past10min"
	PrecipDurPast1H       ParameterID = "precip_dur_past1h"
	PrecipPast10Min       ParameterID = "precip_past10min"
	PrecipPast1H          ParameterID = "precip_past1h"
	PrecipPast1Min        ParameterID = "precip_past1min"
	PrecipPast24H         ParameterID = "precip_past24h"
	Pressure              ParameterID = "pressure"
	PressureAtSea         ParameterID = "pressure_at_sea"
	RadiaGlob             ParameterID = "radia_glob"
	RadiaGlobPast1H       ParameterID = "radia_glob_past1h"
	SnowCoverMan          ParameterID = "snow_cover_man"
	SnowDepthMan          ParameterID = "snow_depth_man"
	SunLast10MinGlob      ParameterID = "sun_last10min_glob"
	SunLast1HGlob         ParameterID = "sun_last1h_glob"
	TempDew               ParameterID = "temp_dew"
	TempDry               ParameterID = "temp_dry"
	TempGrass             ParameterID = "temp_grass"
	TempGrassMaxPast1H    ParameterID = "temp_grass_max_past1h"
	TempGrassMeanPast1H   ParameterID = "temp_grass_mean_past1h"
	TempGrassMinPast1H    ParameterID = "temp_grass_min_past1h"
	TempMaxPast12H        ParameterID = "temp_max_past12h"
	TempMaxPast1H         ParameterID = "temp_max_past1h"
	TempMeanPast1H        ParameterID = "temp_mean_past1h"
	TempMinPast12H        ParameterID = "temp_min_past12h"
	TempMinPast1H         ParameterID = "temp_min_past1h"
	TempSoil              ParameterID = "temp_soil"
	TempSoilMaxPast1H     ParameterID = "temp_soil_max_past1h"
	TempSoilMeanPast1H    ParameterID = "temp_soil_mean_past1h"
	TempSoilMinPast1H     ParameterID = "temp_soil_min_past1h"
	VisibMeanLast10Min    ParameterID = "visib_mean_last10min"
	Visibility            ParameterID = "visibility"
	Weather               ParameterID = "weather"
	WindDir               ParameterID = "wind_dir"
	WindDirPast1H         ParameterID = "wind_dir_past1h"
	WindGustAlwaysPast1H  ParameterID = "wind_gust_always_past1h"
	WindMax               ParameterID = "wind_max"
	WindMaxPer10MinPast1H ParameterID = "wind_max_per10min_past1h"
	WindMin               ParameterID = "wind_min"
	WindMinPast1H         ParameterID = "wind_min_past1h"
	WindSpeed             ParameterID = "wind_speed"
	WindSpeedPast1H       ParameterID = "wind_speed_past1h"
)

type Status string

const (
	Active Status = "Active"
)

type PropertiesType string

const (
	Giws                PropertiesType = "GIWS"
	ManualPrecipitation PropertiesType = "Manual precipitation"
	ManualSnow          PropertiesType = "Manual snow"
	Pluvio              PropertiesType = "Pluvio"
	Synop               PropertiesType = "Synop"
)

type FeatureType string

const (
	Feature FeatureType = "Feature"
)
