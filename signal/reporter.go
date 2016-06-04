package signal

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"
)

// Reporter expresses a generic DC/OS service report
type Reporter interface {
	// Retrieve the endpoints for the service report
	getEndpoints() []string
	// The HTTP method to execute the report retrival
	getMethod() string
	// Retrieve the headers for the HTTP request
	getHeaders() map[string]string
	// Add headers
	addHeaders(map[string]string)
	// Setup the analytics.Track type
	setReport([]byte) error
	// Retreieve the analytics.Track type
	getReport() interface{}
	// Create generic track
	setTrack(config.Config) error
	// Retrieve only track data
	getTrack() *analytics.Track
	// Send track to segmentIO
	sendTrack(config.Config) error
	// Get the name of this Reporter
	getName() string
	// Set an error message
	setError(string)
	// Get an error message
	getError() string
}

// PullReport executes retrival of a service report
func PullReport(r Reporter, c config.Config) error {
	for _, endpoint := range r.getEndpoints() {
		requestURL := fmt.Sprintf("%s%s", c.MasterURL, endpoint)
		log.Debugf("Attempting to pull report from %s", requestURL)
		url, err := url.Parse(requestURL)
		if err != nil {
			return err
		}

		client := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}

		if c.TLSEnabled {
			var tlsClientConfig *tls.Config
			if c.CAPool == nil {
				// do HTTPS without certificate verification.
				tlsClientConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			} else {
				tlsClientConfig = &tls.Config{
					RootCAs: c.CAPool,
				}
			}

			client.Transport = &http.Transport{
				TLSClientConfig: tlsClientConfig,
			}
		}

		req := &http.Request{
			Method: r.getMethod(),
			URL:    url,
			Header: http.Header{},
		}

		headers := r.getHeaders()
		for headerName, headerValue := range headers {
			// ex. headerName = "Content-Type" and headerValue = "application/json"
			req.Header.Add(headerName, headerValue)

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

		if err := r.setReport(body); err != nil {
			return err
		}
	}

	return nil
}
