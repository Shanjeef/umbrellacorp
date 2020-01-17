package weatherforecaster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"umbrellacorp/models"
	"umbrellacorp/util"
)

type openWeatherMap struct {
	baseURL    string
	httpClient http.Client
}

type openWeatherResponse struct {
	List []openWeatherData `json:"list"`
}

// Details from the provider
type openWeatherData struct {
	Dt      int64                `json:"dt"`
	Weather []openWeatherSummary `json:"weather"`
}

func (data openWeatherData) containsWeather(weatherTypes []models.WeatherType) bool {
	for _, weatherType := range weatherTypes {
		openWeatherTypeVal := openWeatherTypeMapping[weatherType]
		for _, summary := range data.Weather {
			if openWeatherTypeVal == summary.Main {
				return true
			}
		}
	}
	return false
}

type openWeatherType string

// models.WeatherType returns models.WeatherType from a openWeatherType value
func (openWeatherKind openWeatherType) WeatherType() (models.WeatherType, bool) {
	for weatherVal, openWeatherVal := range openWeatherTypeMapping {
		if openWeatherKind == openWeatherVal {
			return weatherVal, true
		}
	}
	return models.WeatherTypeRain, false
}

const openWeatherTypeRain = openWeatherType("Rain")

// openWeatherTypeMapping contains mapping of our own standard models.WeatherType to Open Weather's weather types.
var openWeatherTypeMapping = map[models.WeatherType]openWeatherType{
	models.WeatherTypeRain: openWeatherTypeRain,
}

type openWeatherSummary struct {
	Main openWeatherType `json:"main"`
}

func newOpenWeatherMap() *openWeatherMap {
	provider := &openWeatherMap{}
	provider.baseURL = "https://samples.openweathermap.org/data/2.5/forecast"
	provider.httpClient = http.Client{Timeout: 30 * time.Second}
	return provider
}

func (provider *openWeatherMap) UpcomingWeather(city, countrycode string, dateRange util.DateRange, types ...models.WeatherType) ([]models.Weather, error) {
	req, err := http.NewRequest("GET", provider.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request to forecast provider: %s", err.Error())
	}

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("%s,%s", city, countrycode))
	q.Add("appid", "b6907d289e10d714a6e88b30761fae22")
	req.URL.RawQuery = q.Encode()

	resp, err := provider.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request to forecast provider: %v", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response from forecast provider: %v", err.Error())
	}

	defer resp.Body.Close()

	var openWeatherResp openWeatherResponse
	err = json.Unmarshal(body, &openWeatherResp)
	if err != nil {
		return nil, fmt.Errorf("Error parsing response from forecast provider: %v", err.Error())
	}

	return filterAndTranslate(openWeatherResp, dateRange, types)
}

func filterAndTranslate(resp openWeatherResponse, dateRange util.DateRange, types []models.WeatherType) ([]models.Weather, error) {
	weatherDataset := filterData(resp.List, dateRange, types)

	var result []models.Weather
	for _, weatherData := range weatherDataset {
		for _, summary := range weatherData.Weather {
			weatherType, ok := summary.Main.WeatherType()
			if !ok {
				// We don't have a mapping for this weather type, so ignore
				continue
			}

			found := false
			for _, requestedWeatherType := range types {
				if requestedWeatherType == weatherType {
					found = true
					break
				}
			}
			if found {
				result = append(result, models.Weather{
					Date: time.Unix(weatherData.Dt, 0),
					Type: weatherType,
				})
			}
		}

	}
	return result, nil
}

// filterData filters the list of openWeatherData based on the specified dateRange, as well as optionally specified WeatherTypes
func filterData(weatherDataset []openWeatherData, dateRange util.DateRange, types []models.WeatherType) []openWeatherData {
	var result []openWeatherData
	for i, weatherData := range weatherDataset {
		if dateRange.Contains(time.Unix(weatherData.Dt, 0)) {
			if types != nil && weatherData.containsWeather(types) {
				result = append(result, weatherDataset[i])
			}
		}
	}
	return result
}
