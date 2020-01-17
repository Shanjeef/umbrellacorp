package weatherforecaster

import (
	"umbrellacorp/models"
	"umbrellacorp/util"
)

var isMock bool

// Configure configures the pkg in either prod or mock mode
func Configure(mock bool) {
	isMock = mock
}

// Forecaster exposes functionality to retrieve weather details
type Forecaster interface {
	// Obtain upcoming weather for a specific (city, countryCode) combination, with ability to filter for specific weather types within a dateRange.
	// If weather types are not specified, all obtained data from provider is returned
	UpcomingWeather(city, countrycode string, dateRange util.DateRange, types ...models.WeatherType) ([]models.Weather, error)
}

// NewForecaster returns a forecast provider
func NewForecaster() Forecaster {
	if isMock {
		return &mockProvider{}
	}
	return newOpenWeatherMap()
}
