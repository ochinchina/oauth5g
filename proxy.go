package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

// Proxy oauth2 proxy
// Get the token from the authorization server through the proxy
// Verify the token from the authorization server with the key
type Proxy struct {
	router     *gin.Engine
	client     *OAuthClient
	verifier   *AccessTokenVerifier
	tokenCache *TokenCache
}

// NewProxy create a new Proxy object
func NewProxy(tokenReqPath string,
	tokenVerifyPath string,
	oauthServerURL string,
	authServerTLSConfig *tls.Config,
	http2OAuthServer bool,
	tokenVerifyAlgorithm jwa.SignatureAlgorithm,
	key interface{}) *Proxy {
	router := gin.New()
	proxy := &Proxy{router: router,
		client:     NewOAuthClient(oauthServerURL, http2OAuthServer, authServerTLSConfig),
		verifier:   NewAccessTokenVerifier(tokenVerifyAlgorithm, key),
		tokenCache: NewTokenCache(5 * 60)}
	router.POST(tokenReqPath, proxy.HandleTokenRequest)
	router.POST(tokenVerifyPath, proxy.HandleTokenVerify)
	return proxy
}

// Start start proxy, listen on the specified address and accept the token
// access request, the request will be forwarded to the real authorization
// server
func (p *Proxy) Start(addr string) error {
	return p.router.Run(addr)
}

// HandleTokenRequest handle the access token request from the client side
// this request will be forwarded to the authorization server and return
// the access code to the client
func (p *Proxy) HandleTokenRequest(c *gin.Context) {
	atr := NewAccessTokenRequest()
	var err error
	if strings.Contains(c.ContentType(), "application/json") {
		err = atr.FromJSON(c.Request.Body)
	} else {
		err = atr.FromX3WFormEncoding(c.Request.Body)
	}
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
		if resp.FromJSON(r) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		p.cacheTokenFor(atr, resp.ExpiresIn, resp.AccessToken)
		c.Status(http.StatusOK)
		c.Writer.Write(r)
		return
	}
	log.Error("Fail to get the token with error:", err)
	c.Status(http.StatusBadRequest)
}

func (p *Proxy) getTokenFromCache(atr *AccessTokenRequest) (string, error) {
	if atr.IsRequestByType() {
		key := fmt.Sprintf("%s@%s-%s", atr.NfInstanceID, atr.NfType, atr.TargetNfType)
		log.Info("try to get token  by ", key)
		return p.tokenCache.GetToken(key)
	}
	return "", fmt.Errorf("Fail to get token")
}

func (p *Proxy) cacheTokenFor(atr *AccessTokenRequest, expireTime int64, token string) {
	if atr.IsRequestByType() {
		key := fmt.Sprintf("%s@%s-%s", atr.NfInstanceID, atr.NfType, atr.TargetNfType)
		log.Info("Cache the token ", token, " for ", key, " in expire ", expireTime)
		p.tokenCache.CacheToken(key, expireTime, token)
	}
}

// HandleTokenVerify verify the token got from authorization server with the algoritm and the key
func (p *Proxy) HandleTokenVerify(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		log.Error("Fail to read token")
		c.Status(http.StatusBadRequest)
		return
	}
	err = p.verifier.VerifyToken(b)
	if err == nil {
		c.Status(http.StatusOK)
	} else {
		log.Error("Fail to verify token with error:", err)
		c.Status(http.StatusBadRequest)
	}

}

func (p *Proxy) toURLValues(kvs *map[string]interface{}) (url.Values, error) {
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
