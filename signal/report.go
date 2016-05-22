package signal

import (
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

	req := &http.Request{
		Method: r.GetMethod(),
		URL:    url,
	}
	headers := r.GetHeaders()
	if len(headers) > 0 {
		for headerName, headerValue := range headers {
			// ex. headerName = "Content-Type" and headerValue = "application/json"
			req.Header.Set(headerName, headerValue)
		}
	}
	// Add the JWT token to the headers if this is a secure request
	if url.Scheme == "https" {
		if len(c.JWTToken) > 0 {
			bearer := fmt.Sprintf("Bearer %s", c.JWTToken)
			log.Warnf("Authorization: %s", bearer)
			req.Header.Set("Authorization", bearer)
		} else {
			return errors.New("HTTPS requested but no JWT token created.")
		}
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

	if err := r.SetReport(body); err != nil {
		return err
	}

	return nil
}
