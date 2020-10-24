package main

import (
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
)

type AccessTokenClaims struct {
	//  NF instance id of the NRF
	Iss string `json:"iss"`
	// NF instance ID of the NF service consumer
	Sub string `json:"sub"`
	// NF service producer's NF instance ID(s)
	Aud []string `json:"aud"`
	// name of the NF services and the resource/operation-level
	// scopes for which the access_token is authorized for use
	Scope string `json:"scope"`

	// expiration time after which the access_token is considered to be expired
	Exp int64 `json:"exp"`

	// PLMN ID of the NF service consumer
	ConsumerPlmnId *PlmnId `json:"consumerPlmnId,omitempty"`
	// PLMN ID of the NF service producer
	ProducerPlmnId *PlmnId `json:"producerPlmnId,omitempty"`
	// S-NSSAIs of the NF service producer
	ProducerSnssaiList []*Snssai `json:"producerSnssaiList,omitempty"`
	//NSIs of the NF service producer
	ProducerNsiList []string `json:"producerNsiList,omitempty"`
	// NF Set ID of the NF service producer
	ProducerNfSetId string `json:"producerNfSetId,omitempty"`
}

func NewAccessTokenClaims() *AccessTokenClaims {
	return &AccessTokenClaims{}
}

func (atc *AccessTokenClaims) ToJson() ([]byte, error) {
	return json.Marshal(atc)
}
func (atc *AccessTokenClaims) ToJwtToken() jwt.Token {
	token := jwt.New()
	token.Set(jwt.IssuerKey, atc.Iss)
	token.Set(jwt.SubjectKey, atc.Sub)
	token.Set(jwt.AudienceKey, atc.Aud)
	token.Set("scope", atc.Scope)
	token.Set(jwt.ExpirationKey, atc.Exp)
	if atc.ConsumerPlmnId != nil {
		token.Set("consumerPlmnId", atc.ConsumerPlmnId)
	}
	if atc.ProducerPlmnId != nil {
		token.Set("producerPlmnId", atc.ProducerPlmnId)
	}
	if len(atc.ProducerSnssaiList) > 0 {
		token.Set("producerSnssaiList", atc.ProducerSnssaiList)
	}
	if len(atc.ProducerNsiList) > 0 {
		token.Set("producerNsiList", atc.ProducerNsiList)
	}
	if len(atc.ProducerNfSetId) > 0 {
		token.Set("producerNfSetId", atc.ProducerNfSetId)
	}
	return token
}

func (atc *AccessTokenClaims) FromJwtToken(token jwt.Token) error {
	atc.Iss = token.Issuer()
	atc.Sub = token.Subject()
	atc.Aud = token.Audience()
	atc.Exp = token.Expiration().Unix()

	if scope, ok := token.Get("scope"); ok {
		if v, ok := scope.(string); ok {
			atc.Scope = v
		} else {
			return fmt.Errorf("scope must be a string")
		}
	} else {
		return fmt.Errorf("Missing scope field")
	}
	if p, ok := token.Get("consumerPlmnId"); ok {
		if b, err := json.Marshal(p); err == nil {
			atc.ConsumerPlmnId = &PlmnId{}
			if err = json.Unmarshal(b, atc.ConsumerPlmnId); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Fail to decode consumerPlmnId")
		}
	}

	if p, ok := token.Get("producerPlmnId"); ok {
		if b, err := json.Marshal(p); err == nil {
			atc.ProducerPlmnId = &PlmnId{}
			if err = json.Unmarshal(b, atc.ProducerPlmnId); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Fail to decode producerPlmnId")
		}
	}

	if p, ok := token.Get("producerSnssaiList"); ok {
		if b, err := json.Marshal(p); err == nil {
			atc.ProducerSnssaiList = make([]*Snssai, 0)
			if err = json.Unmarshal(b, &atc.ProducerSnssaiList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Fail to decode producerSnssaiList")
		}
	}

	if p, ok := token.Get("producerNsiList"); ok {
		if b, err := json.Marshal(p); err == nil {
			atc.ProducerNsiList = make([]string, 0)
			if err = json.Unmarshal(b, &atc.ProducerNsiList); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Fail to decode producerNsiList")
		}
	}

	if p, ok := token.Get("producerNfSetId"); ok {
		if v, ok := p.(string); ok {
			atc.ProducerNfSetId = v
		} else {
			return fmt.Errorf("producerNfSetId must be string")
		}
	}

	return nil
}
