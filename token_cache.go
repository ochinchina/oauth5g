package main

import (
	"fmt"
	"sync"
	"time"
)

type ExpiryToken struct {
	expireTime int64
	token      string
}
type TokenCache struct {
	sync.Mutex
	tokens map[string]*ExpiryToken
	// minimal left life time for a token in seconds
	minLifeTime int64
}

func NewTokenCache(minLifeTime int64) *TokenCache {
	return &TokenCache{tokens: make(map[string]*ExpiryToken), minLifeTime: minLifeTime}
}

func (stc *TokenCache) CacheToken(key string, expireTime int64, token string) {
	stc.Lock()
	defer stc.Unlock()

	stc.tokens[key] = &ExpiryToken{expireTime: expireTime, token: token}
}

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
