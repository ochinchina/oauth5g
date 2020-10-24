package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"time"
)

type OAuthServer struct {
	// signature algorithm
	alg jwa.SignatureAlgorithm
	// signature key
	key interface{}
	// true to enable http2
	http2       bool
	tlsCertFile string
	tlsKeyFile  string
	router      *gin.Engine
	// the server instance id
	instanceId  string
	tokenExpire time.Duration
	tokenCache  *TokenCache
}

func NewOAuthServer(tokenReqPath string,
	instanceId string,
	tokenExpire time.Duration,
	http2 bool,
	tlsCertFile string,
	tlsKeyFile string,
	alg jwa.SignatureAlgorithm,
	key interface{}) *OAuthServer {
	router := gin.New()
	log.Info("signature algorithm:", alg, ",type of key:", fmt.Sprintf("%T", key))
	server := &OAuthServer{router: router,
		instanceId:  instanceId,
		tokenExpire: tokenExpire,
		http2:       http2,
		tlsCertFile: tlsCertFile,
		tlsKeyFile:  tlsKeyFile,
		alg:         alg,
		key:         key,
		tokenCache:  NewTokenCache(int64(tokenExpire.Seconds() / 2))}
	if len( tokenReqPath ) <= 0 {
		tokenReqPath = "/oauth2/token"
	}
	router.POST(tokenReqPath, server.HandleTokenRequest)
	return server
}

func (s *OAuthServer) Start(addr string) error {
	if s.http2 {
		log.Info("start http2 server")
		h2s := &http2.Server{}
		server := &http.Server{
			Addr:    addr,
			Handler: h2c.NewHandler(s.router, h2s),
		}
		if len(s.tlsCertFile) > 0 && len(s.tlsKeyFile) > 0 {
			return server.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
		} else {
			return server.ListenAndServe()
		}

	} else {
		log.Info("start http server")
		if len(s.tlsCertFile) > 0 && len(s.tlsKeyFile) > 0 {
			return s.router.RunTLS(addr, s.tlsCertFile, s.tlsKeyFile)
		} else {
			return s.router.Run(addr)
		}
	}
}

func (s *OAuthServer) HandleTokenRequest(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		log.Error("Fail to read request with error:", err)
		c.Status(http.StatusBadRequest)
		return
	}
	art := NewAccessTokenRequest()
	err = art.FromX3WFormEncoding(bytes.NewBuffer(b))
	if err != nil {
		log.Error("Fail to decode request with error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	token, err := s.createToken(art)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	resp := AccessTokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.getTokenExpireTime().Unix(),
		Scope:       art.Scope,
	}
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.JSON(http.StatusOK, &resp)
}

func (s *OAuthServer) getTokenFromCache(art *AccessTokenRequest) (string, error) {
	if art.isRequestByType() {
		key := fmt.Sprintf("%s-%s", art.NfType, art.TargetNfType)
		return s.tokenCache.GetToken(key)
	}
	return "", fmt.Errorf("Fail to get token")

}

func (s *OAuthServer) cacheTokenFor(art *AccessTokenRequest, expireTime int64, token string) {
	if art.isRequestByType() {
		key := fmt.Sprintf("%s-%s", art.NfType, art.TargetNfType)
		s.tokenCache.CacheToken(key, expireTime, token)
	}
}

func (s *OAuthServer) createToken(art *AccessTokenRequest) (string, error) {
	b, _ := art.ToJson()
	log.Info("create token from AccessTokenRequest:", string(b))
	if !art.isValid() {
		return "", fmt.Errorf("Not a valid token request")
	}
	t, err := s.getTokenFromCache(art)
	if err == nil {
		return t, nil
	}
	claims, err := s.createClaims(art)
	if err != nil {
		return "", err
	}
	token := claims.ToJwtToken()
	payload, err := jwt.Sign(token, s.alg, s.key)
	if err != nil {
		log.Error("Fail to create JWT Token with error:", err)
		return "", err
	}
	t = string(payload)
	s.cacheTokenFor(art, claims.Exp, t)
	return t, nil
}

func (s *OAuthServer) createClaims(art *AccessTokenRequest) (*AccessTokenClaims, error) {
	if !art.isRequestByType() {
		log.Error("create token only with nfType and targetNfType")
		return nil, fmt.Errorf("Only support create claims by nfType and targetNfType")
	}

	atc := NewAccessTokenClaims()
	atc.Iss = s.instanceId
	atc.Sub = art.NfInstanceId
	if len(art.TargetNfInstanceId) > 0 {
		atc.Aud = []string{art.TargetNfInstanceId}
	} else {
		atc.Aud = []string{art.TargetNfType}
	}

	atc.Scope = art.Scope
	atc.Exp = s.getTokenExpireTime().Unix()
	if art.RequesterPlmn != nil {
		atc.ConsumerPlmnId = art.RequesterPlmn
	}
	if art.TargetPlmn != nil {
		atc.ProducerPlmnId = art.TargetPlmn
	}
	if art.TargetSnssaiList != nil {
		atc.ProducerSnssaiList = art.TargetSnssaiList
	}
	if len(art.TargetNsiList) > 0 {
		atc.ProducerNsiList = art.TargetNsiList
	}
	atc.ProducerNfSetId = art.TargetNfServiceSetId

	return atc, nil
}

func (s *OAuthServer) getTokenExpireTime() time.Time {
	return time.Now().Add(s.tokenExpire)
}
