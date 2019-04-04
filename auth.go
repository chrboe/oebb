package oebb

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const AUTH_URL = "https://tickets.oebb.at/api/domain/v3/init"

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
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
	Cookie             string
}

func auth() (AuthResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", AUTH_URL, nil)
	resp, err := client.Do(req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var authResp AuthResponse
	json.Unmarshal(buf.Bytes(), &authResp)

	return authResp, err
}
