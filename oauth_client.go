package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net"
	"net/http"
)

// OAuthClient used to send AccessTokenRequest to the authorizations server
// to get the access token
type OAuthClient struct {
	serverURL        string
	tlsClientConfig  *tls.Config
	http2OAuthServer bool
}

// NewOAuthClient create a OAuthClient object with:
// - serverURL the authorization server url
// - http2OAuthServer true if the authorization server is a http2 server
// - tlsClientConfig must not be nil if the authorization server is a https server
func NewOAuthClient(serverURL string, http2OAuthServer bool, tlsClientConfig *tls.Config) *OAuthClient {
	return &OAuthClient{serverURL: serverURL,
		http2OAuthServer: http2OAuthServer,
		tlsClientConfig:  tlsClientConfig}
}

// RequestToken request a token from the authorization server
// data - is the AccessTokenRequest object encoded in application/x-www-form-urlencoded
// format
func (oc *OAuthClient) RequestToken(data []byte) ([]byte, error) {
	var client *http.Client = oc.createHTTPClient()

	request, err := http.NewRequest("POST", oc.serverURL, bytes.NewBuffer(data))
	if err != nil {
		log.Error("Fail to request to token from ", oc.serverURL, " with error:", err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if oc.tlsClientConfig != nil && len(oc.tlsClientConfig.ServerName) > 0 {
		request.Host = oc.tlsClientConfig.ServerName
	}

	//resp, err := client.Post(oc.serverURL, "application/x-www-form-urlencoded", bytes.NewBuffer(data) )
	resp, err := client.Do(request)
	if err != nil {
		log.Error("Fail to request to token from ", oc.serverURL, " with error:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		return ioutil.ReadAll(resp.Body)
	}
	log.Error("Fail to get token from ", oc.serverURL, " with status code:", resp.StatusCode)
	return nil, fmt.Errorf("Not 2xx status code %d", resp.StatusCode)
}

func (oc *OAuthClient) createHTTPClient() *http.Client {
	if oc.http2OAuthServer {
		return &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: oc.tlsClientConfig,
				AllowHTTP:       true,
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					if cfg == nil {
						return net.Dial(network, addr)
					}
					return tls.Dial(network, addr, cfg)
				},
			},
		}
	}
	var transport *http.Transport = nil
	if oc.tlsClientConfig != nil {
		transport = &http.Transport{TLSClientConfig: oc.tlsClientConfig}
	}
	return &http.Client{Transport: transport}

}
