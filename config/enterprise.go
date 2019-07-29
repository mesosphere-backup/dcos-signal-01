package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

func initEnterprise() {
	token, err := generateJWTToken()
	if err != nil {
		log.Fatalf("Unable to generate JWT token: %s", err)
		os.Exit(1)
	}

	defaultConfig.DCOSVariant = DCOSVariant{"enterprise"}
	defaultConfig.ExtraHeaders["Authorization"] = fmt.Sprintf("token=%s", token)
}

func generateJWTToken() (string, error) {
	securityConfig := struct {
		UID            string `json:"uid"`
		PrivateKey     string `json:"private_key"`
		LoginEndpoint  string `json:"login_endpoint"`
		JWTToken       string
		SecretJSONPath string
	}{
		SecretJSONPath: "/run/dcos/etc/signal-service/service_account.json",
	}
	// Load the secret file if it exists
	secretJSON, loadErr := ioutil.ReadFile(securityConfig.SecretJSONPath)
	if loadErr != nil {
		log.Warn("Service account not detected, continuing with out secure requests.")
		return "", nil
	}

	if jsonErr := json.Unmarshal(secretJSON, &securityConfig); jsonErr != nil {
		return "", jsonErr
	}

	if securityConfig.UID == "" || securityConfig.PrivateKey == "" || securityConfig.LoginEndpoint == "" {
		return "", errors.New("UID, private key or login endpoint can not be empty.")
	}
	log.Debug("Generating JWT token...")
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["uid"] = securityConfig.UID
	token.Claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenStr, err := token.SignedString([]byte(securityConfig.PrivateKey))
	if err != nil {
		return "", err
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	authReq := struct {
		UID   string `json:"uid"`
		Token string `json:"token,omitempty"`
	}{
		UID:   securityConfig.UID,
		Token: tokenStr,
	}

	b, err := json.Marshal(authReq)
	if err != nil {
		return "", err
	}

	authBody := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", securityConfig.LoginEndpoint, authBody)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to auth with Bouncer,  status code: %d", resp.StatusCode)
	}

	var authResp struct {
		Token string `json:"token"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	log.Debugf("Successfully retrieved JWT token: %s", authResp.Token)
	return authResp.Token, nil
}
