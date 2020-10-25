package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func loadSignatureKeyFromFile(fileName string) (interface{}, error) {
	if strings.HasSuffix(fileName, ".pem") {
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		return loadSignatureKey(b)
	}
	return nil, fmt.Errorf("only .pem file supported")
}
func loadSignatureKey(b []byte) (interface{}, error) {
	p, _ := pem.Decode(b)
	if p == nil {
		return nil, fmt.Errorf("not in PEM format")
	}
	k1, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err == nil {
		return k1, nil
	}
	k2, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err == nil {
		return k2, nil
	}
	k3, err := x509.ParsePKCS1PublicKey(p.Bytes)
	if err == nil {
		return k3, nil
	}
	return x509.ParsePKIXPublicKey(p.Bytes)

}

func toJSONBytes(t interface{}) ([]byte, error) {
	return json.Marshal(t)
}

func toYamlBytes(t interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	encoder := yaml.NewEncoder(buf)
	err := encoder.Encode(t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func loadCertFile(caCertFile string, certFile string, keyFile string) (*tls.Config, error) {
	var cert *tls.Certificate = nil

	if len(certFile) > 0 && len(keyFile) > 0 {
		cert1, err := tls.LoadX509KeyPair(certFile, keyFile)

		if err != nil {
			return nil, err
		}
		cert = &cert1
	}
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		rootCAs = x509.NewCertPool()
	}
	if len(caCertFile) > 0 {
		caCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}

		rootCAs.AppendCertsFromPEM(caCert)
	}

	if cert != nil {
		return &tls.Config{RootCAs: rootCAs, Certificates: []tls.Certificate{*cert}}, nil
	}
	return &tls.Config{RootCAs: rootCAs}, nil
}
