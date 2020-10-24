package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type Proxy struct {
	router     *gin.Engine
	client     *OAuthClient
	verifier   *AccessTokenVerifier
	tokenCache *TokenCache
}

func NewProxy(tokenReqPath string,
	tokenVerifyPath string,
	oauthServerUrl string,
	authServerTlsConfig *tls.Config,
	http2OAuthServer bool,
	tokenVerifyAlgorithm jwa.SignatureAlgorithm,
	key interface{}) *Proxy {
	router := gin.New()
	proxy := &Proxy{router: router,
		client:     NewOAuthClient(oauthServerUrl, http2OAuthServer, authServerTlsConfig),
		verifier:   NewAccessTokenVerifier(tokenVerifyAlgorithm, key),
		tokenCache: NewTokenCache(5 * 60) }
	router.POST(tokenReqPath, proxy.HandleTokenRequest)
	router.POST(tokenVerifyPath, proxy.HandleTokenVerify)
	return proxy
}

func (p *Proxy) Listen(addr string) error {
	return p.router.Run(addr)
}

func (p *Proxy) HandleTokenRequest(c *gin.Context) {
	atr := NewAccessTokenRequest()
	err := atr.FromJson(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	t, err := p.getTokenFromCache(atr)
	if err == nil {
		log.Info("Succeed to get the token:", t, " from local cache")
		c.Status(http.StatusOK)
		c.Writer.Write([]byte(t))
		return
	}
	b, err := atr.ToX3WFormEncoding()

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if r, err := p.client.RequestToken(b); err == nil {
		log.Info("Succeed to get the token:", string(r), " from remote server")
		resp := NewAccessTokenResponse()
		if resp.FromBytes(r) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		p.cacheTokenFor(atr, resp.ExpiresIn, resp.AccessToken)
		c.Status(http.StatusOK)
		c.Writer.Write(r)
		return
	} else {
		log.Error("Fail to get the token with error:", err)
		c.Status(http.StatusBadRequest)
	}
}

func (p *Proxy) getTokenFromCache(atr *AccessTokenRequest) (string, error) {
	if atr.isRequestByType() {
		key := fmt.Sprintf("%s-%s", atr.NfType, atr.TargetNfType)
		log.Info("try to get token  by ", key)
		return p.tokenCache.GetToken(key)
	}
	return "", fmt.Errorf("Fail to get token")
}

func (p *Proxy) cacheTokenFor(atr *AccessTokenRequest, expireTime int64, token string) {
	if atr.isRequestByType() {
		key := fmt.Sprintf("%s-%s", atr.NfType, atr.TargetNfType)
		log.Info("Cache the token ", token, " for ", key, " in expire ", expireTime)
		p.tokenCache.CacheToken(key, expireTime, token)
	}
}

func (p *Proxy) HandleTokenVerify(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		log.Error("Fail to read token")
		c.Status(http.StatusBadRequest)
		return
	}
	err = p.verifier.VerifyToken( b )
	if err == nil {
		c.Status(http.StatusOK)
	} else {
		log.Error("Fail to verify token with error:", err)
		c.Status(http.StatusBadRequest)
	}

}

func (p *Proxy) toUrlValues(kvs *map[string]interface{}) (url.Values, error) {
	values := url.Values{}

	for k, v := range *kvs {
		if s, ok := v.(string); ok {
			values.Add(k, s)
		} else {
			return url.Values{}, fmt.Errorf("%s value is not a string", k)
		}
	}
	return values, nil
}
