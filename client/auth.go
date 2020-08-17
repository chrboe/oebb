package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const authURL = "https://tickets.oebb.at/api/domain/v3/init"

type authResponse struct {
	AccessToken string `json:"accessToken"` // this is actually duplicated in the response
	Token       struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"token"`
	Channel            string    `json:"channel"`
	SupportID          string    `json:"supportId"`
	CashID             string    `json:"cashId"`
	OrgUnit            int       `json:"orgUnit"`
	LegacyUserMigrated bool      `json:"legacyUserMigrated"`
	UserID             string    `json:"userId"`
	PersonID           string    `json:"personId"`
	CustomerID         string    `json:"customerId"`
	Realm              string    `json:"realm"`
	SessionID          string    `json:"sessionId"`
	SessionTimeout     int       `json:"sessionTimeout"`
	SessionVersion     string    `json:"sessionVersion"`
	SessionCreatedAt   time.Time `json:"sessionCreatedAt"`
	XffxIP             string    `json:"xffxIP"`
}

// AuthInfo describes info used to authenticate requests against the ÖBB API.
// This information can be cached and re-used at a later time.
type AuthInfo struct {
	AccessToken string
	Channel     string
	SessionID   string
	SupportID   string
	ExpiresIn   int
}

// Auth authenticates against the ÖBB API.
// The result is a response containing, among other things, an access token and
// a refresh token.
// No credentials are actually required to interact with the API.
func Auth() (AuthInfo, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", authURL, nil)
	resp, err := client.Do(req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var authResp authResponse
	json.Unmarshal(buf.Bytes(), &authResp)

	info := AuthInfo{
		AccessToken: authResp.Token.AccessToken,
		Channel:     authResp.Channel,
		SessionID:   authResp.SessionID,
		SupportID:   authResp.SupportID,
		ExpiresIn:   authResp.SessionTimeout,
	}

	return info, err
}
