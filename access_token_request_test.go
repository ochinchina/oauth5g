package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAccessTokenRequestMandatoryFields(t *testing.T) {
	//s := `{"grant_type":"client_credentials","nfInstanceId":"123","scope":"NMF","requesterPlmnList":[{"mcc":"281","123"}]}`
	atr := NewAccessTokenRequest()
	atr.GrantType = "client_credentials"
	atr.NfInstanceId = "123"
	atr.Scope = "NMF"
	p1 := PlmnId{Mcc: "081", Mnc: "123"}
	atr.RequesterPlmnList = []*PlmnId{&p1}
	b, _ := atr.ToJson()
	s := string(b)
	err := atr.FromJson(bytes.NewBufferString(s))
	if err != nil {
		t.Fail()
	}
	b, err = atr.ToJson()
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(b))
	s1, err := atr.ToX3WFormEncoding()
	if err != nil {
		t.Fail()
	}
	fmt.Println(s1)
}

func TestAccessTokenRequestDecodeFromEncoding(t *testing.T) {
	s := "grant_type=client_credentials&nfInstanceId=123&scope=NMF&requesterPlmnList.0.mcc=281&requesterPlmnList.0.mnc=123"
	atr := NewAccessTokenRequest()
	r := bytes.NewBufferString(s)
	err := atr.FromX3WFormEncoding(r)
	if err != nil {
		t.Fail()
	}
	b, _ := atr.ToJson()
	fmt.Println(string(b))
}
