package main

import (
	"encoding/json"
	"io"
)

// AccessTokenResponse the access token response from the
// authorization server to the client.
// A AccessTokenResponse object will be replied to client
// after granting the access to specific producer
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// NewAccessTokenResponse create a new AccessTokenResponse object
func NewAccessTokenResponse() *AccessTokenResponse {
	return &AccessTokenResponse{}
}

// FromJSON create a AccessTokenResponse from json format
func (atr *AccessTokenResponse) FromJSON(b []byte) error {
	return json.Unmarshal(b, atr)
}

// FromReader create a AccessTokenResponse from a json reader
func (atr *AccessTokenResponse) FromReader(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(atr)
}

// ToJSON convert the AccessTokenResponse object to json format
func (atr *AccessTokenResponse) ToJSON() ([]byte, error) {
	return json.Marshal(atr)
}
