package main

import (
	"encoding/json"
	"io"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func NewAccessTokenResponse() *AccessTokenResponse {
	return &AccessTokenResponse{}
}

func (atr *AccessTokenResponse) FromBytes(b []byte) error {
	return json.Unmarshal(b, atr)
}

func (atr *AccessTokenResponse) FromReader(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(atr)
}

func (atr *AccessTokenResponse) ToJson() ([]byte, error) {
	return json.Marshal(atr)
}
