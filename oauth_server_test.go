package main

import (
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"testing"
	"time"
)

var privateKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAln8DNKCaqc1PeksxHdKpa4bO8bq+bQFdcJwV0dHYFKzXk7Vk
9TeSHeikmFD5QDSZ6Lb5r3d+iE3Zj/g4vmrwttaynzt+PdTNh6bWAS+J2S8JePeC
MC+d4A55XtAZdnn/wzOFPivQVE4ny09QsZfLKJcJFRqVgIWQ+Qx3G905wSdCTyjh
zSmxXio3HXrsWwut9RLae6c/oUAQibhJZ61RpFKcJK/+9Yp0PIIDykkTAJQ48yry
MunK4BPAHyGwJMICq0ruinMihzl5wQ+K+b1hUWmglO/lk7TUQlTWzwTZhM8tbRSm
QbJjWcT95JAP+e9TMwUc0UzvJqQoJSZwUtcmCwIDAQABAoIBABRy87ON8xYJgLgk
SWs8+apTqgDsl/1lxxWLD+fvtIWyqqQ2bJ5pS6BmwR61OfkAbL1TR4ARI49PzLUk
TEmLbGPbrm+2fMPYC4lYIRIOLJPnMUlPkFAN2Ezeip1Yd21CzH7wnDVDbM5XxE56
24MWFhPQ2foVH6VUAMvmZrTRjYxyHNVhvCeO1CkPRpbyuMZ+UsbJKDd/xsR5IzpU
gkOw6ZaAuVbX8HOx7fBC6sIRqvo6ZBqnpPAmauciPbKKekI4hpWKEjS0Y4ILJjfs
+0RX833x8192JlbLenXnAQQxPa03fNdh1AkDiN8Jkp5LU/lRkuQZTYPcDzCklwvr
Xct79VECgYEAyCkfS84+HhyfQtN7rZAc6dyzwDuvjtCa68DqPy8NFK4Cqnsifo+v
zy7qDHeI+EuFOuT3I4cooU9+cqPU+z9YYwz9YZGmkW7p71bWwzz9wFaRmlWMpscI
LA44FVJqwWgBrHkHDVBCk8I6Q7cjT6bDg9tnJ9CyoRfIe8iTACER/xMCgYEAwHr/
9/0qxWyfk5p7cbL0P8P4NrhdKcAPbulpmz+ipRqVSdgz22yZT7HzEU/OTfBiNUAH
mtDM49XVL3JEjEi2ydjfT1Bn0/qYYVxcDGC/v6DK71o5IqgrfEQN1oZIVsUPOEDv
Okfk4Cm9ra0JyxY1uNA1eDP8yzTIa87WQa9DBCkCgYEAnUrbgim/3M/nQ4+HuHvw
jMLYKq98pYE+zRcbva1O3Tplc+0xzT1DDlAysrtY0q4eM4rsv2mePy2GE7a1Tv+X
iLcTgxH/UHhVs7SNLn4GdphQ8XRbBFCSFnTSE8dhhz0hW5T8OrUgrJbMTJxlTlmh
eUP4S2yQg1F6RfP3uPlD+CsCgYEAhxeNAdOZKlkzotgg3csY7IwxcM5y7LOU4WZH
LaQ7FiATOXHZ655L+AhQLg1SIZeehfs7mygDNcFFz/gmLkN2rzJcgQFQ7hGK04KM
RE+/JNLIu7caNL3NT3lAMRmsOeIy7Wt9u+zrsXz6WKQDJJug9uaDMKtkOIcCR9Ay
xoUoxwkCgYEAj7UjmKqgaiO0Gw2rTCarDkV8MjLanY84e16P7vaX7sTqTcjPV7rK
cE4JGbrFDX/MOKest1t8DLHt126uho8AkTXv8DSynJ73tT+FK7pO3JFyQwhnldO4
uwXXeJThGJPsnLrmwyLTXtGgdoUmD6QB+gY0XgiMDyu03ErDL3SthnI=
-----END RSA PRIVATE KEY-----`

func createToken() (string, error) {
	key, err := loadSignatureKey([]byte(privateKey))
	if err != nil {
		return "", err
	}
	server := NewOAuthServer("", "instance-1", time.Duration(3600)*time.Second, false, "", "", jwa.RS256, key)
	req := NewAccessTokenRequest()
	req.GrantType = "client_credentials"
	req.NfInstanceID = "12345"
	req.NfType = "LMF"
	req.TargetNfType = "AMF"
	req.Scope = "namf-comm"
	req.TargetNsiList = []string{"nsi-1", "nsi-2", "nsi-3"}
	req.TargetSnssaiList = []*Snssai{&Snssai{Sst: 10}, &Snssai{Sst: 30, Sd: "sd-2"}}
	token, accessTokenErr := server.createToken(req)
	if accessTokenErr != nil {
		return "", fmt.Errorf(accessTokenErr.Error)
	}
	return token, nil

}

func TestClaims(t *testing.T) {
	_, err := createToken()
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
}
