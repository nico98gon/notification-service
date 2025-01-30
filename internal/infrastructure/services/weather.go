package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Dia     string `json:"Dia"`
		Tempo   string `json:"Tempo"`
		Maxima  string `json:"Maxima"`
		Minima  string `json:"Minima"`
		IUV     string `json:"IUV"`
	} `json:"data"`
}

type WaveResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Period        string  `json:"Period"`
		Agitation     string  `json:"Agitation"`
		Height        float64 `json:"Height"`
		Direction     string  `json:"Direction"`
		WindSpeed     float64 `json:"WindSpeed"`
		WindDirection string  `json:"WindDirection"`
	} `json:"data"`
}

func GetWeatherForecast(cityID int, isCoastal bool) (string, string, error) {
	weatherURL := fmt.Sprintf("http://weather-service:8083/city-forecast?city_id=%d", cityID)
	resp, err := http.Get(weatherURL)
	if err != nil {
		return "", "", fmt.Errorf("error al conectar con weather-service: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error al leer respuesta del weather-service: %v", err)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return "", "", fmt.Errorf("error al parsear respuesta JSON: %v", err)
	}

	if len(weatherData.Data) < 4 {
		return "", "", fmt.Errorf("weather-service no devolvió suficientes datos")
	}

	forecastContent := ""
	for i := 0; i < 4; i++ {
		forecast := weatherData.Data[i]
		forecastContent += fmt.Sprintf("%s: %s, Máx: %s°C, Mín: %s°C, IUV: %s\n", forecast.Dia, forecast.Tempo, forecast.Maxima, forecast.Minima, forecast.IUV)
	}

	waveContent := ""
	if isCoastal {
		waveURL := fmt.Sprintf("http://weather-service:8083/wave-forecast?city_id=%d&day=0", cityID)
		resp, err := http.Get(waveURL)
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				var waveData WaveResponse
				if err := json.Unmarshal(body, &waveData); err == nil && waveData.Status == "success" && len(waveData.Data) > 0 {
					waveContent = "\nPrevisión de olas:\n"
					for _, wave := range waveData.Data {
						waveContent += fmt.Sprintf("%s - Agitación: %s, Altura: %.1fm, Dirección: %s, Viento: %.1fm/s (%s)\n",
							wave.Period, wave.Agitation, wave.Height, wave.Direction, wave.WindSpeed, wave.WindDirection)
					}
				}
			}
		}
	}

	title := "Pronóstico del tiempo"
	content := forecastContent + waveContent
	return title, content, nil
}
