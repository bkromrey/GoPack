// Package weatherAPI utilizes the Geocoding API and Open-Meteo API to look up
// the forecast data for a specified location passed in as a string.
//
// If the forecast data can be found, the resulting structure will contain the
// forecast for a range of dates contained in a map, and if the forecast data
// was not successful, then this will return an error.
//
// The Geocoding API key must be placed in .env in the top level directory.
//
// Open-Meteo has the following restrictions on the free API:
// * forecast data only available 16 days into the future and 3 months in past.

package weatherAPI

import (
	"GoPack/env"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// A ForecastData contains forecast data (high & low temperatures) as well as
// the chance of precipitation for a single day.
type ForecastData struct {
	TempMax             float64
	TempMin             float64
	PrecipitationChance float64
}

// A WeatherResults contains one of two things, either a map where each entry
// represents a single day's forecast data or an error message if the Open-Meteo
// API was not able to return the requested data.
type WeatherResults struct {
	Data     map[string]ForecastData `json:"Data"`
	ErrorMsg string                  `json:"ErrorMsg"`
}

// getLatLong looks up the lat & long coordinates of a location defined by a string
func getLatLong(location string) (string, string) {

	url := fmt.Sprintf("https://geocode.maps.co/search?q=%v&api_key=%v", location, env.GetGeoAPIKey())

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, sendErr := http.DefaultClient.Do(request)
	if sendErr != nil {
		log.Fatal(err)
	}

	// close response when we're done with it
	defer response.Body.Close()

	// process request
	locationData, openErr := io.ReadAll(response.Body)
	if openErr != nil {
		log.Fatal(err)
	}

	var parsedLocationData []map[string]interface{}
	unpackErr := json.Unmarshal(locationData, &parsedLocationData)
	if unpackErr != nil {
		log.Fatal(unpackErr)
	}

	data := parsedLocationData[0]

	latitude := data["lat"].(string)
	longitude := data["lon"].(string)

	return latitude, longitude
}

// The callWeatherAPI makes a call to the Open-Meteo API to retrieve weather forecast
// data from Open-Meteo.
//
// Requires the latitude & longitude be passed in as strings, and the start and
// end dates must be formatted in YYYY-MM-DD format.
func callWeatherAPI(lat string, long string, startDate string, endDate string) map[string]interface{} {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%v&longitude=%v&daily=temperature_2m_max,temperature_2m_min,precipitation_probability_max&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch&timezone=auto&start_date=%v&end_date=%v", lat, long, startDate, endDate)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, sendErr := http.DefaultClient.Do(request)
	if sendErr != nil {
		log.Fatal(sendErr)
	}

	// close the response when we're done with it
	defer response.Body.Close()

	// process body of the request
	forecastJSON, openErr := io.ReadAll(response.Body)
	if openErr != nil {
		log.Fatal(err)
	}

	var parsedData map[string]interface{}
	unpackErr := json.Unmarshal(forecastJSON, &parsedData)
	if unpackErr != nil {
		log.Fatal(unpackErr)
	}

	return parsedData
}

func processForecastJSON(rawData map[string]interface{}) map[string]ForecastData {

	// first we need to type cast our raw data
	dailyData, ok := rawData["daily"].(map[string]interface{})
	if !ok {
		log.Fatal("there was an issue typecasting")
	}

	timeData, ok := dailyData["time"].([]interface{})
	if !ok {
		log.Fatal("there was an issue typecasting timeData")
	}

	maxTempData, ok := dailyData["temperature_2m_max"].([]interface{})
	if !ok {
		log.Fatal("there was an issue typecasting maxTempData")
	}

	minTempData, ok := dailyData["temperature_2m_min"].([]interface{})
	if !ok {
		log.Fatal("there was an issue typecasting minTempData")
	}

	precipitationData, ok := dailyData["precipitation_probability_max"].([]interface{})
	if !ok {
		log.Fatal("there was an issue typecasting precipitationData")
	}

	// then organize data into usable structure
	processedData := make(map[string]ForecastData)

	index := 0
	for index < len(dailyData) {

		processedData[timeData[index].(string)] = ForecastData{
			TempMax:             maxTempData[index].(float64),
			TempMin:             minTempData[index].(float64),
			PrecipitationChance: precipitationData[index].(float64),
		}
		index += 1
	}
	return processedData
}

// The GetWeatherFromAPI function is the main function another package calls to
// get the forecast data for a range of dates based on a location name.
//
// The dates must be formatted in a string in YYYY-MM-DD format and the location
// should also be a string.
func GetWeatherFromAPI(location string, departDate string, returnDate string) WeatherResults {

	// geocode the location string
	lat, long := getLatLong(location)

	// request forecast & precipitation data
	forecastDataRaw := callWeatherAPI(lat, long, departDate, returnDate)

	// parse data into usable structure
	var results WeatherResults

	if forecastDataRaw["error"] != true {
		forecast := processForecastJSON(forecastDataRaw)
		results.Data = forecast

	} else {
		results.ErrorMsg = forecastDataRaw["reason"].(string)
	}

	return results

}
