package main

import (
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
)

// AccessTokenClaims the token claims defined in the 5G
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
	ConsumerPlmnID *PlmnID `json:"consumerPlmnID,omitempty"`
	// PLMN ID of the NF service producer
	ProducerPlmnID *PlmnID `json:"producerPlmnID,omitempty"`
	// S-NSSAIs of the NF service producer
	ProducerSnssaiList []*Snssai `json:"producerSnssaiList,omitempty"`
	//NSIs of the NF service producer
	ProducerNsiList []string `json:"producerNsiList,omitempty"`
	// NF Set ID of the NF service producer
	ProducerNfSetID string `json:"producerNfSetId,omitempty"`
}

// NewAccessTokenClaims create a AccessTokenClaims object
func NewAccessTokenClaims() *AccessTokenClaims {
	return &AccessTokenClaims{}
}

// ToJSON convert AccessTokenClaims object to json format
func (atc *AccessTokenClaims) ToJSON() ([]byte, error) {
	return json.Marshal(atc)
}

// ToJwtToken convert the AccessTokenClaims object to jwt.Token object
func (atc *AccessTokenClaims) ToJwtToken() jwt.Token {
	token := jwt.New()
	token.Set(jwt.IssuerKey, atc.Iss)
	token.Set(jwt.SubjectKey, atc.Sub)
	token.Set(jwt.AudienceKey, atc.Aud)
	token.Set("scope", atc.Scope)
	token.Set(jwt.ExpirationKey, atc.Exp)
	if atc.ConsumerPlmnID != nil {
		token.Set("consumerPlmnID", atc.ConsumerPlmnID)
	}
	if atc.ProducerPlmnID != nil {
		token.Set("producerPlmnID", atc.ProducerPlmnID)
	}
	if len(atc.ProducerSnssaiList) > 0 {
		token.Set("producerSnssaiList", atc.ProducerSnssaiList)
	}
	if len(atc.ProducerNsiList) > 0 {
		token.Set("producerNsiList", atc.ProducerNsiList)
	}
	if len(atc.ProducerNfSetID) > 0 {
		token.Set("producerNfSetId", atc.ProducerNfSetID)
	}
	return token
}

// FromJwtToken create AccessTokenClaims from a jwt.Token
// if some mandatory fields are missing, return error
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
	if p, ok := token.Get("consumerPlmnID"); ok {
		if b, err := json.Marshal(p); err == nil {
			atc.ConsumerPlmnID = &PlmnID{}
			if err = json.Unmarshal(b, atc.ConsumerPlmnID); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Fail to decode consumerPlmnID")
		}
	}

	if p, ok := token.Get("producerPlmnID"); ok {
		if b, err := json.Marshal(p); err == nil {
			atc.ProducerPlmnID = &PlmnID{}
			if err = json.Unmarshal(b, atc.ProducerPlmnID); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Fail to decode producerPlmnID")
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
			atc.ProducerNfSetID = v
		} else {
			return fmt.Errorf("producerNfSetId must be string")
		}
	}

	return nil
}
