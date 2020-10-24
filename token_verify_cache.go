package main

import (
	"sync"
	"time"
)

// TokenVerifyCache caches the verified tokens
//
type TokenVerifyCache struct {
	sync.Mutex
	tokens map[string]int64
}

func NewTokenVerifyCache() *TokenVerifyCache {
	return &TokenVerifyCache{ tokens: make( map[string]int64 ) }
}


// AddVerifiedToken add a verified token with its expire time to the cache
//
// if the expireTime is less than current time, the token will not be added
// to the cache
func (tvc *TokenVerifyCache) AddVerifiedToken( token string, expireTime int64 ) {
	if expireTime <= time.Now().Unix() {
		return
	}

	tvc.Lock()
	defer tvc.Unlock()

	tvc.tokens[token] = expireTime
}

// IsTokenVerified check if the token is verified and it is not expired
//
// return true if the token is verifed and not expired
func (tvc *TokenVerifyCache) IsTokenVerified( token string ) bool {
	tvc.Lock()
        defer tvc.Unlock()

	if v, ok := tvc.tokens[ token ]; ok {
		if v > time.Now().Unix() {
			return true
		}
		delete( tvc.tokens, token )
	}
	return false
}
