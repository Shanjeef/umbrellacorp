package weatherforecaster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"umbrellacorp/models"
	"umbrellacorp/util"
)

type mockProvider struct{}

func (provider *mockProvider) UpcomingWeather(city, countrycode string, dateRange util.DateRange, types ...models.WeatherType) ([]models.Weather, error) {
	buf, err := ioutil.ReadFile("mock_response.json")
	if err != nil {
		return nil, fmt.Errorf("Couldn't read mock response file: %s", err.Error())
	}

	var resp openWeatherResponse
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal response details: %s", err.Error())
	}

	return filterAndTranslate(resp, dateRange, types)
}
