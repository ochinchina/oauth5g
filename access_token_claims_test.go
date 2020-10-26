package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"testing"
)

func TestFromToJwtToken(t *testing.T) {
	atc := NewAccessTokenClaims()

	atc.Iss = "NRF"
	atc.Sub = "test-124"
	atc.Exp = 1603178453
	atc.Aud = []string{"AMF"}
	atc.Scope = "namf-comm"
	atc.ProducerSnssaiList = []*Snssai{&Snssai{Sst: 10}, &Snssai{Sst: 20, Sd: "sd-1"}}
	atc.ConsumerPlmnID = &PlmnID{Mcc: "123", Mnc: "456"}
	atc.ProducerPlmnID = &PlmnID{Mcc: "123", Mnc: "789"}
	atc.ProducerNsiList = []string{"nsi-1", "nsi-2"}
	atc.ProducerNfSetID = "nf-set-id"
	token := atc.ToJwtToken()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fail()
	}

	payload, err := jwt.Sign(token, jwa.RS256, privateKey)
	if err != nil {
		t.Fail()
	}
	newToken, err := jwt.ParseBytes(payload)

	fmt.Printf("newToken:%v\n", newToken)
	newAtc := NewAccessTokenClaims()
	err = newAtc.FromJwtToken(newToken)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		t.Fail()
	}
	fmt.Printf("new claims:%v\n", newAtc)
	b, _ := newAtc.ToJSON()
	fmt.Printf("%s\n", string(b))
}
