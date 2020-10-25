package main

import (
	"bytes"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"time"
)

// AccessTokenVerifier defined to verify if the token
// is a valid token. When a producer received a request
// from consumer, it will extract the access token from
// the "Authorization" header and verify the token with
// the specified signature algorithm and the key
type AccessTokenVerifier struct {
	alg                jwa.SignatureAlgorithm
	key                interface{}
	verifiedTokenCache *TokenVerifyCache
}

// NewAccessTokenVerifier create a AccessTokenVerifier object with the specific signature
// algorithm and its key
func NewAccessTokenVerifier(alg jwa.SignatureAlgorithm, key interface{}) *AccessTokenVerifier {
	log.Info("create verifier with algorithm ", alg, " and key ", fmt.Sprintf("%T", key))
	return &AccessTokenVerifier{alg: alg,
		key:                key,
		verifiedTokenCache: NewTokenVerifyCache()}
}

// VerifyToken verify the token with the signature algoritm and the key. If the token
// is valid and not expired, return nil
func (atv *AccessTokenVerifier) VerifyToken(b []byte) error {

	if atv.verifiedTokenCache.IsTokenVerified(string(b)) {
		return nil
	}

	if atv.key == nil {
		return fmt.Errorf("Fail to verify token because key is nil")
	}
	token, err := jwt.Parse(bytes.NewBuffer(b), jwt.WithVerify(atv.alg, atv.key))
	if err != nil {
		return err
	}
	atc := NewAccessTokenClaims()
	err = atc.FromJwtToken(token)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if atc.Exp < now {
		return fmt.Errorf("Expiration time %d is less than current time %d", atc.Exp, now)
	} else {
		atv.verifiedTokenCache.AddVerifiedToken(string(b), atc.Exp)
	}

	return err
}
