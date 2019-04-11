package oebb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const connectionsURL = "https://tickets.oebb.at/api/hafas/v4/timetable"

//
// REQUEST
//

type connectionRequest struct {
	Reverse           bool                   `json:"reverse"`
	DatetimeDeparture string                 `json:"datetimeDeparture"`
	Filter            connectionsFilter      `json:"filter"`
	Passengers        []passenger            `json:"passengers"`
	Count             int                    `json:"count"`
	DebugFilter       connectionsDebugFilter `json:"debugFilter"`
	SortType          string                 `json:"sortType"`
	From              Station                `json:"from"`
	To                Station                `json:"to"`
	Timeout           struct{}               `json:"timeout"`
}

type connectionsFilter struct {
	Regionaltrains     bool `json:"regionaltrains"`
	Direct             bool `json:"direct"`
	ChangeTime         bool `json:"changeTime"`
	Wheelchair         bool `json:"wheelchair"`
	Bikes              bool `json:"bikes"`
	Trains             bool `json:"trains"`
	Motorail           bool `json:"motorail"`
	DroppedConnections bool `json:"droppedConnections"`
}

type challengedFlags struct {
	HasHandicappedPass bool `json:"hasHandicappedPass"`
	HasAssistanceDog   bool `json:"hasAssistanceDog"`
	HasWheelchair      bool `json:"hasWheelchair"`
	HasAttendant       bool `json:"hasAttendant"`
}

type passenger struct {
	Type                string          `json:"type"`
	ID                  int             `json:"id"`
	Me                  bool            `json:"me"`
	Remembered          bool            `json:"remembered"`
	ChallengedFlags     challengedFlags `json:"challengedFlags"`
	Relations           []interface{}   `json:"relations"`
	Cards               []interface{}   `json:"cards"`
	BirthdateChangeable bool            `json:"birthdateChangeable"`
	BirthdateDeletable  bool            `json:"birthdateDeletable"`
	NameChangeable      bool            `json:"nameChangeable"`
	PassengerDeletable  bool            `json:"passengerDeletable"`
	IsSelected          bool            `json:"isSelected"`
}

type connectionsDebugFilter struct {
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

type connectionsResponse struct {
	Connections []Connection `json:"connections"`
}

type DepartureStation struct {
	Name                       string `json:"name"`
	Esn                        int    `json:"esn"`
	Departure                  string `json:"departure"`
	DepartureDelay             string `json:"departureDelay"`
	DeparturePlatform          string `json:"departurePlatform"`
	DeparturePlatformDeviation string `json:"departurePlatformDeviation"`
	ShowAsResolvedMetaStation  bool   `json:"showAsResolvedMetaStation"`
}

type ArrivalStation struct {
	Name                      string `json:"name"`
	Esn                       int    `json:"esn"`
	Arrival                   string `json:"arrival"`
	ArrivalDelay              string `json:"arrivalDelay"`
	ArrivalPlatform           string `json:"arrivalPlatform"`
	ArrivalPlatformDeviation  string `json:"arrivalPlatformDeviation"`
	ShowAsResolvedMetaStation bool   `json:"showAsResolvedMetaStation"`
}

type TranslatedString struct {
	De string `json:"de"`
	En string `json:"en"`
	It string `json:"it"`
}

type Category struct {
	Name                          string           `json:"name"`
	Number                        string           `json:"number"`
	ShortName                     string           `json:"shortName"`
	DisplayName                   string           `json:"displayName"`
	LongName                      TranslatedString `json:"longName"`
	BackgroundColor               string           `json:"backgroundColor"`
	FontColor                     string           `json:"fontColor"`
	BarColor                      string           `json:"barColor"`
	Place                         TranslatedString `json:"place"`
	JourneyPreviewIconID          string           `json:"journeyPreviewIconId"`
	JourneyPreviewIconColor       string           `json:"journeyPreviewIconColor"`
	AssistantIconID               string           `json:"assistantIconId"`
	Train                         bool             `json:"train"`
	ParallelLongName              string           `json:"parallelLongName"`
	ParallelDisplayName           string           `json:"parallelDisplayName"`
	BackgroundColorDisabledMobile string           `json:"backgroundColorDisabledMobile"`
	BackgroundColorDisabled       string           `json:"backgroundColorDisabled"`
	FontColorDisabled             string           `json:"fontColorDisabled"`
	BarColorDisabled              string           `json:"barColorDisabled"`
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

func GetConnections(from, to Station, a AuthInfo, departureTime time.Time, numResults int) ([]Connection, error) {
	client := &http.Client{}

	cr := connectionRequest{
		Reverse:           false,
		DatetimeDeparture: departureTime.Format("2006-01-02T15:04:05.999"),
		Filter: connectionsFilter{
			Regionaltrains:     false,
			Direct:             false,
			ChangeTime:         false,
			Wheelchair:         false,
			Bikes:              false,
			Trains:             false,
			Motorail:           false,
			DroppedConnections: false,
		},
		Passengers: []passenger{
			passenger{
				Type:       "ADULT",
				ID:         1554277150,
				Me:         false,
				Remembered: false,
				ChallengedFlags: challengedFlags{
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
		Count: numResults,
		DebugFilter: connectionsDebugFilter{
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

	req, err := http.NewRequest("POST", connectionsURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Channel", a.Channel)
	req.Header.Add("AccessToken", a.AccessToken)
	req.Header.Add("SessionId", a.SessionID)
	req.Header.Add("x-ts-supportid", "WEB_"+a.SupportID)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	connections := &connectionsResponse{}
	json.Unmarshal(buf.Bytes(), connections)

	return connections.Connections, err
}
