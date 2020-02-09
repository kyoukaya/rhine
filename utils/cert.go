package utils

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/elazarl/goproxy"
)

// LoadCA loads the cert and RSA key pair from the specified paths and configures
// goproxy to use it.
func LoadCA(certPath, keyPath string) error {
	var err error
	certFile, err := os.Open(BinDir + certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	keyFile, err := os.Open(BinDir + keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	caCert, err := ioutil.ReadAll(certFile)
	if err != nil {
		return err
	}
	caKey, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return err
	}
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}
