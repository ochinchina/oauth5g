package main

import (
	"fmt"
	"testing"
)

func TestAccessTokenResponseDecode(t *testing.T) {
	s := `{
       "access_token":"2YotnFZFEjr1zCsicMWpAA",
       "token_type":"example",
       "expires_in":3600,
       "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
       "example_parameter":"example_value"
     }`

	r := NewAccessTokenResponse()

	err := r.FromBytes([]byte(s))
	if err != nil {
		t.Fail()
	}
	if r.AccessToken != "2YotnFZFEjr1zCsicMWpAA" {
		t.Fail()
	}

	if r.ExpiresIn != 3600 {
		t.Fail()
	}
	b, err := r.ToJson()
	if err != nil {
		t.Fail()
	}

	fmt.Println(string(b))
}
