package oebb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const CONNECTIONS_URL = "https://tickets.oebb.at/api/hafas/v4/timetable"

//
// REQUEST
//

type ConnectionsRequest struct {
	Reverse           bool        `json:"reverse"`
	DatetimeDeparture string      `json:"datetimeDeparture"`
	Filter            Filter      `json:"filter"`
	Passengers        []Passenger `json:"passengers"`
	Count             int         `json:"count"`
	DebugFilter       DebugFilter `json:"debugFilter"`
	SortType          string      `json:"sortType"`
	From              Station     `json:"from"`
	To                Station     `json:"to"`
	Timeout           struct{}    `json:"timeout"`
}

type Filter struct {
	Regionaltrains     bool `json:"regionaltrains"`
	Direct             bool `json:"direct"`
	ChangeTime         bool `json:"changeTime"`
	Wheelchair         bool `json:"wheelchair"`
	Bikes              bool `json:"bikes"`
	Trains             bool `json:"trains"`
	Motorail           bool `json:"motorail"`
	DroppedConnections bool `json:"droppedConnections"`
}

type ChallengedFlags struct {
	HasHandicappedPass bool `json:"hasHandicappedPass"`
	HasAssistanceDog   bool `json:"hasAssistanceDog"`
	HasWheelchair      bool `json:"hasWheelchair"`
	HasAttendant       bool `json:"hasAttendant"`
}

type Passenger struct {
	Type                string          `json:"type"`
	ID                  int             `json:"id"`
	Me                  bool            `json:"me"`
	Remembered          bool            `json:"remembered"`
	ChallengedFlags     ChallengedFlags `json:"challengedFlags"`
	Relations           []interface{}   `json:"relations"`
	Cards               []interface{}   `json:"cards"`
	BirthdateChangeable bool            `json:"birthdateChangeable"`
	BirthdateDeletable  bool            `json:"birthdateDeletable"`
	NameChangeable      bool            `json:"nameChangeable"`
	PassengerDeletable  bool            `json:"passengerDeletable"`
	IsSelected          bool            `json:"isSelected"`
}

type DebugFilter struct {
	NoAggregationFilter bool `json:"noAggregationFilter"`
	NoEqclassFilter     bool `json:"noEqclassFilter"`
	NoNrtpathFilter     bool `json:"noNrtpathFilter"`
	NoPaymentFilter     bool `json:"noPaymentFilter"`
	UseTripartFilter    bool `json:"useTripartFilter"`
	NoVbxFilter         bool `json:"noVbxFilter"`
	NoCategoriesFilter  bool `json:"noCategoriesFilter"`
}

//
// RESPONSE
//

type ConnectionsResponse struct {
	Connections []Connection `json:"connections"`
}

type DepartureStation struct {
	Name                       string `json:"name"`
	Esn                        int    `json:"esn"`
	Departure                  string `json:"departure"`
	DeparturePlatform          string `json:"departurePlatform"`
	DeparturePlatformDeviation string `json:"departurePlatformDeviation"`
	ShowAsResolvedMetaStation  bool   `json:"showAsResolvedMetaStation"`
}

type ArrivalStation struct {
	Name                      string `json:"name"`
	Esn                       int    `json:"esn"`
	Arrival                   string `json:"arrival"`
	ArrivalPlatform           string `json:"arrivalPlatform"`
	ArrivalPlatformDeviation  string `json:"arrivalPlatform"`
	ShowAsResolvedMetaStation bool   `json:"showAsResolvedMetaStation"`
}

type LongName struct {
	De string `json:"de"`
	En string `json:"en"`
	It string `json:"it"`
}

type Place struct {
	De string `json:"de"`
	En string `json:"en"`
	It string `json:"it"`
}

type Category struct {
	Name                          string   `json:"name"`
	Number                        string   `json:"number"`
	ShortName                     string   `json:"shortName"`
	DisplayName                   string   `json:"displayName"`
	LongName                      LongName `json:"longName"`
	BackgroundColor               string   `json:"backgroundColor"`
	FontColor                     string   `json:"fontColor"`
	BarColor                      string   `json:"barColor"`
	Place                         Place    `json:"place"`
	JourneyPreviewIconID          string   `json:"journeyPreviewIconId"`
	JourneyPreviewIconColor       string   `json:"journeyPreviewIconColor"`
	AssistantIconID               string   `json:"assistantIconId"`
	Train                         bool     `json:"train"`
	ParallelLongName              string   `json:"parallelLongName"`
	ParallelDisplayName           string   `json:"parallelDisplayName"`
	BackgroundColorDisabledMobile string   `json:"backgroundColorDisabledMobile"`
	BackgroundColorDisabled       string   `json:"backgroundColorDisabled"`
	FontColorDisabled             string   `json:"fontColorDisabled"`
	BarColorDisabled              string   `json:"barColorDisabled"`
}

type Section struct {
	From        DepartureStation `json:"from,omitempty"`
	To          ArrivalStation   `json:"to,omitempty"`
	Duration    int              `json:"duration"`
	Category    Category         `json:"category,omitempty"`
	Type        string           `json:"type"`
	HasRealtime bool             `json:"hasRealtime"`
}

type Connection struct {
	ID       string           `json:"id"`
	From     DepartureStation `json:"from"`
	To       ArrivalStation   `json:"to"`
	Sections []Section        `json:"sections"`
	Switches int              `json:"switches"`
	Duration int              `json:"duration"`
}

func getConnections(from, to Station, a AuthResponse) (*ConnectionsResponse, error) {
	client := &http.Client{}

	cr := ConnectionsRequest{
		Reverse:           false,
		DatetimeDeparture: time.Now().Format("2006-01-02T15:04:05.999"),
		Filter: Filter{
			Regionaltrains:     false,
			Direct:             false,
			ChangeTime:         false,
			Wheelchair:         false,
			Bikes:              false,
			Trains:             false,
			Motorail:           false,
			DroppedConnections: false,
		},
		Passengers: []Passenger{
			Passenger{
				Type:       "ADULT",
				ID:         1554277150,
				Me:         false,
				Remembered: false,
				ChallengedFlags: ChallengedFlags{
					HasHandicappedPass: false,
					HasAssistanceDog:   false,
					HasWheelchair:      false,
					HasAttendant:       false,
				},
				Relations:           []interface{}{},
				Cards:               []interface{}{},
				BirthdateChangeable: true,
				BirthdateDeletable:  true,
				NameChangeable:      true,
				PassengerDeletable:  true,
				IsSelected:          false,
			},
		},
		Count: 5,
		DebugFilter: DebugFilter{
			NoAggregationFilter: false,
			NoEqclassFilter:     false,
			NoNrtpathFilter:     false,
			NoPaymentFilter:     false,
			UseTripartFilter:    false,
			NoVbxFilter:         false,
			NoCategoriesFilter:  false,
		},
		SortType: "DEPARTURE",
		From:     from,
		To:       to,
	}

	body, err := json.Marshal(cr)

	req, err := http.NewRequest("POST", CONNECTIONS_URL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Channel", a.Channel)
	req.Header.Add("AccessToken", a.AccessToken)
	req.Header.Add("SessionId", a.SessionID)
	req.Header.Add("x-ts-supportid", "WEB_"+a.SupportID)

	req.AddCookie(&http.Cookie{Name: "ts-cookie", Value: a.Cookie})
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	connections := &ConnectionsResponse{}
	json.Unmarshal(buf.Bytes(), connections)

	return connections, err
}
