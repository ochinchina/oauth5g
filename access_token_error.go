package main

import (
	"bytes"
	"encoding/json"
)

const (
	InvalidRequest       string = "invalid_request"
	InvalidClient        string = "invalid_client"
	InvalidGrant         string = "invalid_grant"
	UnauthorizedClient   string = "unauthorized_client"
	UnsupportedGrantType string = "unsupported_grant_type"
	InvalidScope         string = "invalid_scope"
)

// AccessTokenError the error return to oauth client if error happens
//
// Except set the status code to non-200, the server should return
// a AccessTokenError json object to the client
type AccessTokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}

// NewAccessTokenError create a AccessTokenError object
func NewAccessTokenError(err string) *AccessTokenError {
	return &AccessTokenError{Error: err}
}

// ToJSON convert the AccessTokenError to json object
func (ate *AccessTokenError) ToJSON() ([]byte, error) {
	return json.Marshal(ate)
}

// FromJSON create a AccessTokenError from a json
func (ate *AccessTokenError) FromJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	return decoder.Decode(ate)
}
