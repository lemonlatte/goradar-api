package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GORADAR_URI = "stop_fucking_with_us.goradar.io"
	GORADAR_UA  = "GoRadar/13 CFNetwork/711.1.16 Darwin/14.0.0"
)

type GoRadarData struct {
	Pokemons []PokemonLocation `json:"pokemons"`
}

type PokemonLocation struct {
	Id            string  `json:"encounter_id"`
	PokemonId     int64   `json:"pokemon_id"`
	Name          string  `json:"pokemon_name"`
	DisappearTime int64   `json:"disappear_time"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
}

func GetPokemon(request func(*http.Request) (*http.Response, error), swLat, swLng, neLat, neLng float64) (v *GoRadarData, err error) {
	queryUrl := url.URL{
		Scheme: "https",
		Host:   GORADAR_URI,
		Path:   "/raw_data",
		RawQuery: fmt.Sprintf(
			"swLat=%f&swLng=%f&neLat=%f&neLng=%f&pokemon=false&pokestops=false&gyms=false",
			swLat, swLng, neLat, neLng,
		),
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", GORADAR_UA)

	resp, err := request(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		v = &GoRadarData{}
		d := json.NewDecoder(resp.Body)
		err = d.Decode(v)
		if err != nil {
			return
		}
	}
	return
}
