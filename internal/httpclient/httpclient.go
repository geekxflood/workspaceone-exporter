// internal/httpclient/httpclient.go

package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

func New(insecure bool) *http.Client {
	tlsConfig := &tls.Config{InsecureSkipVerify: insecure}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Adjust timeout as needed
	}
}
