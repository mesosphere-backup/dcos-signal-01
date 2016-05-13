package signal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

func pullHealthReport(healthURL string, endpoint string) (hr *HealthReport, err error) {
	log.Debug("Attempting to pull health report from ", healthURL, endpoint)

	url := fmt.Sprintf("%s%s", healthURL, endpoint)
	client := http.Client{
		Timeout: time.Duration(time.Second),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error handling response from ", url, ": ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := json.Unmarshal(body, &hr); err != nil {
		log.Error("Error unmarshaling JSON body: ", err)
		return nil, err
	}

	return hr, nil
}
