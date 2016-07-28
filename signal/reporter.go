package signal

import (
	"bytes"
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
	appendError(string)
	// Get an error message
	getError() []string
}

// PullReport executes retrival of a service report
func PullReport(endpoint string, r Reporter, c config.Config) error {
	url, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	log.Debugf("Pulling from %s", endpoint)
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	if url.Scheme == "https" {
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

	urlStr := fmt.Sprintf("%v", url)
	method := r.getMethod()
	reqBody := "{}"
	req, _ := http.NewRequest(method, urlStr, bytes.NewBufferString(reqBody))

	headers := r.getHeaders()
	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}
	log.Debugf("Request %s: %+v", endpoint, req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response %s %s: %s", resp.Proto, endpoint, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("Response %s: %s, proto %s", resp.Proto, endpoint, resp.Status)

	if err := r.setReport(body); err != nil {
		return err
	}

	return nil
}
