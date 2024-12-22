package tui

import (
	"GoPack/fileHandling"
	"GoPack/weatherAPI"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"sort"
)

type WeatherScreenModel struct {
	list     fileHandling.PackingList
	forecast map[string]weatherAPI.ForecastData
	errorMsg string
	display  string
}

func InitWeatherScreenModel(list fileHandling.PackingList) WeatherScreenModel {
	return WeatherScreenModel{list: list, display: "Loading..."}
}

func (m WeatherScreenModel) Init() tea.Cmd {
	return nil
}

func (m WeatherScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case forecastMsg:
		if msg.forecastResponse.ErrorMsg == "" {
			m.forecast = msg.forecastResponse.Data
			return m, formatForecast(m)
		} else {
			m.errorMsg = msg.forecastResponse.ErrorMsg
			return m, formatError(m)
		}

	case formattedForecastMsg:
		m.display = msg.data

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		}

		// press any key to return to list
		return PackingListTUI, GoToPackingListObj(m.list)

	}

	return m, nil

}

func (m WeatherScreenModel) View() string {
	return m.display
}

// COMMANDS -------------------------------------------------------------
func getWeather(list fileHandling.PackingList) tea.Cmd {
	return func() tea.Msg {

		// request forecast data
		//forecastResponse := zmqClient.SendWeatherRequest(list.Destination, list.DepartDate, list.ReturnDate)
		forecastResponse := weatherAPI.GetWeatherFromAPI(list.Destination, list.DepartDate, list.ReturnDate)

		return forecastMsg{forecastResponse: forecastResponse}
	}
}

type forecastMsg struct {
	forecastResponse weatherAPI.WeatherResults
}

func formatForecast(m WeatherScreenModel) tea.Cmd {
	return func() tea.Msg {

		header := fmt.Sprintf("Forecast for %v from %v to %v\n\n",
			m.list.Destination, m.list.DepartDate, m.list.ReturnDate)
		header += fmt.Sprintf("%17v%6v%21v\n", "HIGH", "LOW", "CHANCE OF RAIN/SNOW")

		// sort & display data for each date from oldest date to newest date
		var dates []string
		for date := range m.forecast {
			dates = append(dates, date)
		}

		sort.Strings(dates)

		// then put all the forecast data together
		var main string
		for _, date := range dates {
			main += fmt.Sprintf("%v    %2.0f°   %2.0f°                  %2.0f%%\n", date, m.forecast[date].TempMax, m.forecast[date].TempMin, m.forecast[date].PrecipitationChance)
		}

		footer := "\n\nPress any key to return to packing list."
		return formattedForecastMsg{
			data: header + main + footer,
		}
	}
}

func formatError(m WeatherScreenModel) tea.Cmd {
	return func() tea.Msg {
		return formattedForecastMsg{data: m.errorMsg}
	}
}

type formattedForecastMsg struct {
	data string
}
