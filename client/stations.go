package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const stationsURL = "https://tickets.oebb.at/api/hafas/v1/stations"

type Station struct {
	Latitude  int `json:"latitude"`
	Longitude int `json:"longitude"`
	// Name in this context is actually either the "name" or the "meta"
	// field of the station. Meta stands for a station name which is not
	// actually a single station, but rather a whole city or a group of
	// stations. If the name is not set in the response, the "meta" field
	// is used instead. If the name is set, "meta" is ignored (but it is
	// usually empty anyways).
	Name string `json:"name"`
	// Along with the "name or meta" field, we also store here whether or
	// not this is a meta station (but we don't transfer that information
	// to JSON, because the API wouldn't be happy with it).
	Meta   bool `json:"-"`
	Number int  `json:"number"`
}

func (s *Station) UnmarshalJSON(data []byte) error {
	var raw struct {
		Latitude  int    `json:"latitude"`
		Longitude int    `json:"longitude"`
		Name      string `json:"name"`
		Meta      string `json:"meta"`
		Number    int    `json:"number"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.Latitude = raw.Latitude
	s.Longitude = raw.Longitude
	s.Name = raw.Name
	if s.Name == "" {
		s.Name = raw.Meta
	}
	s.Meta = raw.Meta != ""
	s.Number = raw.Number

	return nil
}

func GetStations(name string, a AuthInfo) ([]Station, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", stationsURL+"?name="+url.QueryEscape(name), nil)
	req.Header.Add("Channel", a.Channel)
	req.Header.Add("AccessToken", a.AccessToken)
	req.Header.Add("SessionId", a.SessionID)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 440:
		// login time-out: session expird
		return nil, &SessionTimeoutError{}
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var stations []Station
	json.Unmarshal(buf.Bytes(), &stations)

	return stations, nil
}
