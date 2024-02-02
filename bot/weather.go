package bot

import (
	"context"

	owm "github.com/briandowns/openweathermap"
)

// WeatherService represents a service for retrieving weather data.
type WeatherService struct {
	context.Context
}

// NewWeatherClient creates a new WeatherService instance.
func NewWeatherClient(ctx context.Context) WeatherService {
	return WeatherService{
		Context: ctx,
	}
}

// Weather is an interface for weather-related operations.
type Weather interface {
	GetWeatherByLocation(string) (string, error)
}

// GetWeatherByLocation retrieves current weather data for a specified location.
func (s WeatherService) GetWeatherByLocation(req string) (*owm.CurrentWeatherData, error) {
	// Call the internal GetWeatherByLocation function to retrieve weather data.
	res, err := GetWeatherByLocation(req) // removed handler.GetWeatherByLocation(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetWeatherByLocation retrieves current weather data from OpenWeatherMap for a specified location.
func GetWeatherByLocation(location string) (*owm.CurrentWeatherData, error) {
	// Create a new instance of OpenWeatherMap's CurrentWeatherData with English output.
	w, err := owm.NewCurrent("C", "EN", OpenWeatherToken) // (internal - OpenWeatherMap reference for kelvin) with English output
	if err != nil {
		return nil, err
	}

	// Retrieve current weather data for the specified location.
	if err := w.CurrentByName(location); err != nil {
		return nil, err
	}

	return w, nil
}
