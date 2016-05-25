package signal

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

type Reporter interface {
	SetEndpoints([]string)
	GetEndpoints() []string
	SetMethod(string)
	GetMethod() string
	SetHeaders(map[string]string)
	GetHeaders() map[string]string
	SetReport([]byte) error
	GetReport() interface{}
	SetTrack(config.Config) error
	GetTrack() *analytics.Track
	SendTrack(config.Config) error
}

func PullReport(r Reporter, c config.Config) error {
	for _, endpoint := range r.GetEndpoints() {
		requestURL := fmt.Sprintf("%s%s", c.MasterURL, endpoint)
		log.Debugf("Attempting to pull report from %s", requestURL)
		url, err := url.Parse(requestURL)
		if err != nil {
			return err
		}

		client := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}

		if url.Scheme == "https" {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		} else if url.Scheme != "http" {
			return errors.New(fmt.Sprintf("Transport protocol not supported: %s", url.Scheme))
		}

		req := &http.Request{
			Method: r.GetMethod(),
			URL:    url,
			Header: http.Header{},
		}
		headers := r.GetHeaders()
		for headerName, headerValue := range headers {
			// ex. headerName = "Content-Type" and headerValue = "application/json"
			req.Header.Add(headerName, headerValue)

		}
		// Add the JWT token to the headers if this is a secure request
		if len(c.JWTToken) > 0 {
			bearer := fmt.Sprintf("token=%s", c.JWTToken)
			// Removing this for production, here for debugging
			log.Debug("Making request with authorization bearer")
			req.Header.Set("Authorization", bearer)
		} else {
			log.Warn("No JWT token present, making anonymous request")
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Debugf("Successful request to %s", requestURL)

		if err := r.SetReport(body); err != nil {
			return err
		}
	}

	return nil
}
