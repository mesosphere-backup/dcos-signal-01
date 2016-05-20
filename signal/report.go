package signal

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dcos/dcos-signal/config"
)

type Report interface {
	SetURL(string)
	GetURL() string
	SetMethod(string)
	GetMethod() string
	SetHeaders(map[string]string)
	GetHeaders() map[string]string
	SetReport([]byte) error
	GetReport() interface{}
	Track(config.Config) error
}

func PullReport(r Report, c config.Config) error {
	log.Warn("TLS Disabled")
	url := r.GetURL()
	method := r.GetMethod()
	headers := r.GetHeaders()

	log.Debugf("Attempting to pull report from %s", url)
	client := http.Client{
		Timeout: time.Duration(time.Second),
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for headerName, headerValue := range headers {
			// ex. headerName = "Content-Type" and headerValue = "application/json"
			req.Header.Set(headerName, headerValue)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := r.SetReport(body); err != nil {
		return err
	}

	return nil
}
