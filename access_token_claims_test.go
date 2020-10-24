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
	atc.ConsumerPlmnId = &PlmnId{Mcc: "123", Mnc: "456"}
	atc.ProducerPlmnId = &PlmnId{Mcc: "123", Mnc: "789"}
	atc.ProducerNsiList = []string{"nsi-1", "nsi-2"}
	atc.ProducerNfSetId = "nf-set-id"
	token := atc.ToJwtToken()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fail()
	}

	payload, err := jwt.Sign(token, jwa.RS256, privateKey)
	if err != nil {
		t.Fail()
	}
	new_token, err := jwt.ParseBytes(payload)

	fmt.Printf("new_token:%v\n", new_token)
	new_atc := NewAccessTokenClaims()
	err = new_atc.FromJwtToken(new_token)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		t.Fail()
	}
	fmt.Printf("new claims:%v\n", new_atc)
	b, _ := new_atc.ToJson()
	fmt.Printf("%s\n", string(b))
}
