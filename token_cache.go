package main

import (
	"fmt"
	"sync"
	"time"
)

// ExpiryToken a access token with expiry time
type ExpiryToken struct {
	expireTime int64
	token      string
}

// TokenCache the created token cache.
// The authorization server will cache its created acess token
// for reuse if a new similar request is within the minLifeTime
type TokenCache struct {
	sync.Mutex
	tokens map[string]*ExpiryToken
	// minimal left life time for a token in seconds
	minLifeTime int64
}

// NewTokenCache create a TokenCache with minLifeTime in seconds
func NewTokenCache(minLifeTime int64) *TokenCache {
	return &TokenCache{tokens: make(map[string]*ExpiryToken), minLifeTime: minLifeTime}
}

// CacheToken cache the token with the expireTime
func (stc *TokenCache) CacheToken(key string, expireTime int64, token string) {
	stc.Lock()
	defer stc.Unlock()

	stc.tokens[key] = &ExpiryToken{expireTime: expireTime, token: token}
}

// GetToken get token by key. A valid key will be return if the token of the
// key exists and the token is not expired within minLifeTime
func (stc *TokenCache) GetToken(key string) (string, error) {
	stc.Lock()
	defer stc.Unlock()

	stc.clearExpiredTokens()

	if v, ok := stc.tokens[key]; ok && v.expireTime > time.Now().Unix()+stc.minLifeTime {
		return v.token, nil
	}
	return "", fmt.Errorf("No token for %s", key)
}

func (stc *TokenCache) clearExpiredTokens() {
	expiredKeys := make([]string, 0)
	for k, v := range stc.tokens {
		if v.expireTime < time.Now().Unix() {
			expiredKeys = append(expiredKeys, k)
		}
	}
	for _, k := range expiredKeys {
		delete(stc.tokens, k)
	}
}
