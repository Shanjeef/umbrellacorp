package weatherforecaster

import (
	"testing"
	"time"
	"umbrellacorp/models"
	"umbrellacorp/util"

	"github.com/stretchr/testify/assert"
)

func TestOpenWeather(t *testing.T) {
	Configure(true)

	now := time.Date(2017, 02, 16, 0, 0, 0, 0, time.UTC)
	dateRange := util.DateRange{Start: now, End: now.AddDate(0, 0, 5)}
	weatherDetails, err := NewForecaster().UpcomingWeather("Toronto", "CA", dateRange, models.WeatherTypeRain)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 11, len(weatherDetails))
}
