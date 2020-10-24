package main

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"testing"
)

var public_key string = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAln8DNKCaqc1PeksxHdKp
a4bO8bq+bQFdcJwV0dHYFKzXk7Vk9TeSHeikmFD5QDSZ6Lb5r3d+iE3Zj/g4vmrw
ttaynzt+PdTNh6bWAS+J2S8JePeCMC+d4A55XtAZdnn/wzOFPivQVE4ny09QsZfL
KJcJFRqVgIWQ+Qx3G905wSdCTyjhzSmxXio3HXrsWwut9RLae6c/oUAQibhJZ61R
pFKcJK/+9Yp0PIIDykkTAJQ48yryMunK4BPAHyGwJMICq0ruinMihzl5wQ+K+b1h
UWmglO/lk7TUQlTWzwTZhM8tbRSmQbJjWcT95JAP+e9TMwUc0UzvJqQoJSZwUtcm
CwIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	key, err := loadSignatureKey([]byte(public_key))
	if err != nil {
		t.Fail()
	}

	token, _ := createToken()


	verifier := NewAccessTokenVerifier(jwa.RS256, key)
	err = verifier.VerifyToken([]byte(token))
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
}
