package oebb

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const STATIONS_URL = "https://tickets.oebb.at/api/hafas/v1/stations"

type Station struct {
	Latitude  int    `json:"latitude"`
	Longitude int    `json:"longitude"`
	Name      string `json:"name"`
	Meta      string `json:"meta"`
	Number    int    `json:"number"`
}

func getStations(name string, a AuthResponse) ([]Station, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", STATIONS_URL+"?name="+url.QueryEscape(name), nil)
	req.Header.Add("Channel", a.Channel)
	req.Header.Add("AccessToken", a.AccessToken)
	req.Header.Add("SessionId", a.SessionID)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var stations []Station
	json.Unmarshal(buf.Bytes(), &stations)

	return stations, nil
}
