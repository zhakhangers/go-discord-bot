package bot

import (
	"context"

	owm "github.com/briandowns/openweathermap"
)

type WeatherService struct {
	context.Context
}

func NewWeatherClient(ctx context.Context) WeatherService {
	return WeatherService{
		Context: ctx,
	}
}

type Weather interface {
	GetWeatherByLocation(string) (string, error)
}

func (s WeatherService) GetWeatherByLocation(req string) (*owm.CurrentWeatherData, error) {
	res, err := GetWeatherByLocation(req) // removed handler.GetWeatherByLocation(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// const (
// 	openWeatherApiKey = "73c1bf24555f8bf022059c86d63a7fad"
// )

func GetWeatherByLocation(location string) (*owm.CurrentWeatherData, error) {
	w, err := owm.NewCurrent("C", "EN", OpenWeatherToken) // (internal - OpenWeatherMap reference for kelvin) with English output
	if err != nil {
		return nil, err
	}

	if err := w.CurrentByName(location); err != nil {
		return nil, err
	}

	return w, nil
}
