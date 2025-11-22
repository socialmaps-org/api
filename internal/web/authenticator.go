package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Authentication struct {
	Active           bool   `json:"active"`
	ClientID         string `json:"client_id"`
	Exp              uint64 `json:"exp"`
	Iat              uint64 `json:"iat"`
	Scope            string `json:"scope"`
	Username         string `json:"username"`
	OpenStreetMapSub string `json:"openstreetmap_sub"`
}

type Authenticator interface {
	Introspect(token, scope string) (*Authentication, error)
}

type authenticator struct {
	url          string
	clientID     string
	clientSecret string
}

func NewAuthenticator(url, clientID, clientSecret string) Authenticator {
	return &authenticator{
		url:          url,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (a *authenticator) Introspect(token, scope string) (*Authentication, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		a.url,
		strings.NewReader(
			url.Values{
				"token":           []string{token},
				"scope":           []string{scope},
				"token_type_hint": []string{"access_token"},
			}.Encode(),
		),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(a.clientID, a.clientSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection failure: got HTTP %d", res.StatusCode)
	}

	var doc Authentication
	err = json.NewDecoder(res.Body).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
