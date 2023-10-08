package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/weather", getWeather)

	err := http.ListenAndServe(":3333", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is my website!\n")
}
func getWeather(w http.ResponseWriter, r *http.Request) {
	citiesQueryString := r.URL.Query().Get("cities")
	cities := strings.Split(citiesQueryString, ",")
	results := []*weather{}
	errs := []error{}

	for _, c := range cities {
		weather, err := getWeatherInCity(c)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, weather)
	}

	if len(errs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		json.NewEncoder(w).Encode(errs)
        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(results)
}

type location struct {
	Name string `json:"name"`
}

type current struct {
	Temperature float32 `json:"temp_c"`
}

type weatherApiResponse struct {
	Location location `json:"location"`
	Current  current  `json:"current"`
}

type weather struct {
	City        string  `json:"city"`
	Temperature float32 `json:"temperature"`
}

func getWeatherInCity(city string) (w *weather, err error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if len(apiKey) == 0 {
        println("apikeyerr")
		return nil, errors.New("Api key not found")
	}
	requestUrl := getRequestUrl(city, apiKey)
	response, err := http.Get(requestUrl)

	if err != nil {
		return nil, err
	}

	var weatherApiResponseJson weatherApiResponse
	err = json.NewDecoder(response.Body).Decode(&weatherApiResponseJson)

	if err != nil {
		return nil, err
	}

	w = &weather{
		City:        weatherApiResponseJson.Location.Name,
		Temperature: weatherApiResponseJson.Current.Temperature,
	}

	return w, err
}

func getRequestUrl(city string, apiKey string) string {
	return fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=yes", apiKey, city)
}
