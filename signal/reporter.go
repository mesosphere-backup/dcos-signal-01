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
	SetURL(string)
	GetURL() string
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
	log.Debugf("Attempting to pull report from %s", r.GetURL())
	url, err := url.Parse(r.GetURL())
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: time.Duration(time.Second),
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
	if len(headers) > 0 {
		for headerName, headerValue := range headers {
			// ex. headerName = "Content-Type" and headerValue = "application/json"
			req.Header.Add(headerName, headerValue)

		}
	}
	// Add the JWT token to the headers if this is a secure request
	if len(c.JWTToken) > 0 {
		bearer := fmt.Sprintf("token=%s", c.JWTToken)
		// Removing this for production, here for debugging
		log.Warnf("HTTPS Enabled: Authorization: %s", bearer)
		req.Header.Set("Authorization", bearer)
	} else {
		log.Warn("No JWT token present, making insecure request.")
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("Successful request to %s", r.GetURL())

	if err := r.SetReport(body); err != nil {
		return err
	}

	return nil
}
