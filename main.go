package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	GORADAR_URI = "stop_fucking_with_us.goradar.io"
	GORADAR_UA  = "GoRadar/13 CFNetwork/711.1.16 Darwin/14.0.0"
)

type GoRadarData struct {
	Pokemons []PokemonLocation `json:"pokemons"`
}

type PokemonLocation struct {
	Id        int     `json:"pokemon_id"`
	Name      string  `json:"pokemon_name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func getPokemon(request func(*http.Request) (*http.Response, error), swLat, swLng, neLat, neLng float64) (v *GoRadarData, err error) {
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

func main() {
	c := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			MaxIdleConnsPerHost:   20,
			TLSHandshakeTimeout:   3 * time.Second,
			ResponseHeaderTimeout: 3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	d, err := getPokemon(c.Do, 25.043603, 121.558737, 25.076540, 121.581137)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", d)
}
