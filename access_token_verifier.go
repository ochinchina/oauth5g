package main

import (
	"bytes"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"time"
)

type AccessTokenVerifier struct {
	alg jwa.SignatureAlgorithm
	key interface{}
	verifiedTokenCache *TokenVerifyCache
}

func NewAccessTokenVerifier(alg jwa.SignatureAlgorithm, key interface{}) *AccessTokenVerifier {
	log.Info("create verifier with algorithm ", alg, " and key ", fmt.Sprintf("%T", key))
	return &AccessTokenVerifier{alg: alg,
				key: key,
				verifiedTokenCache: NewTokenVerifyCache() }
}

func (atv *AccessTokenVerifier) VerifyToken(b []byte) error {

	if atv.verifiedTokenCache.IsTokenVerified(string(b) ) {
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
		atv.verifiedTokenCache.AddVerifiedToken( string(b), atc.Exp )
	}

	return err
}
