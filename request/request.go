package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestData struct {
	User   string `json:"User"`
	Secret string `json:"Secret,omitempty"`
}

func ConstructPostRequest(baseURL string, port int, endpoint string, data *RequestData) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%d/%s", baseURL, port, endpoint), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func ExecuteRequest(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func PostSignup(client *http.Client, baseURL string, port int, user string) ([]byte, error) {
	data := &RequestData{User: user}
	req, err := ConstructPostRequest(baseURL, port, "signup", data)
	if err != nil {
		return nil, err
	}
	return ExecuteRequest(client, req)
}

func PostCheck(client *http.Client, baseURL string, port int, user string) ([]byte, error) {
	data := &RequestData{User: user}
	req, err := ConstructPostRequest(baseURL, port, "check", data)
	if err != nil {
		return nil, err
	}
	return ExecuteRequest(client, req)
}

func PostGetUserSecret(client *http.Client, baseURL string, port int, user string) (string, error) {
	data := &RequestData{User: user}
	req, err := ConstructPostRequest(baseURL, port, "getUserSecret", data)
	if err != nil {
		return "", err
	}
	respBody, err := ExecuteRequest(client, req)
	if err != nil {
		return "", err
	}

	const prefix = "User secret: "
	if strings.HasPrefix(string(respBody), prefix) {
		return strings.TrimSpace(strings.TrimPrefix(string(respBody), prefix)), nil
	}
	return "", fmt.Errorf("Unexpected format for getUserSecretBody")
}

func PostGetUserLevel(client *http.Client, baseURL string, port int, user, secret string) ([]byte, error) {
	data := &RequestData{User: user, Secret: secret}
	req, err := ConstructPostRequest(baseURL, port, "getUserLevel", data)
	if err != nil {
		return nil, err
	}
	return ExecuteRequest(client, req)
}

func PostGetUserPoints(client *http.Client, baseURL string, port int, user, secret string) ([]byte, error) {
	data := &RequestData{User: user, Secret: secret}
	req, err := ConstructPostRequest(baseURL, port, "getUserPoints", data)
	if err != nil {
		return nil, err
	}
	return ExecuteRequest(client, req)
}
