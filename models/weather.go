package models

import "time"

// Weather details
type Weather struct {
	Date time.Time   `json:"date"`
	Type WeatherType `json:"type"`
}

// WeatherType outlines type of weather
type WeatherType string

// WeatherTypeRain signifies rain
const WeatherTypeRain = WeatherType("Rain")
