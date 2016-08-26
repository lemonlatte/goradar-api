package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/lemonlatte/goradar-api/api"
)

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

	d, err := api.GetPokemon(c.Do, 25.043603, 121.558737, 25.076540, 121.581137)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", d)
}
