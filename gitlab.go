package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type gitlab struct {
	host  string
	token string
}

type deployKey struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

func gitlabClient(host string, token string) (*gitlab, error) {
	if host == "" {
		return nil, fmt.Errorf("host string cannot be empty or null")
	}

	if token == "" {
		return nil, fmt.Errorf("token string cannot be empty or null")
	}

	return &gitlab{host: host, token: token}, nil
}

func unmarshalError(body []byte, statusCode int) error {
	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		return err
	}
	return fmt.Errorf("http error %d: %s", statusCode, dat)
}

func (g *gitlab) deployKey(projectID string, key string, push bool) (*deployKey, error) {
	url := fmt.Sprintf("https://%s/api/v4/projects/%s/deploy_keys", g.host, projectID)
	timestamp := time.Now().Format("20060102150405")
	bodyRaw := `{
		"title": "vault-%s",
		"key": "%s",
		"can_push": %t
	}`
	body := fmt.Sprintf(bodyRaw, timestamp, key, push)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", g.token)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, unmarshalError(respBody, resp.StatusCode)
	}

	deployKey := deployKey{}

	if err := json.Unmarshal(respBody, &deployKey); err != nil {
		return nil, err
	}

	return &deployKey, nil
}

func (g *gitlab) delDeployKey(projectID string, keyID string, b *backend) (bool, error) {
	url := fmt.Sprintf("https://%s/api/v4/projects/%s/deploy_keys/%s", g.host, projectID, keyID)
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("PRIVATE-TOKEN", g.token)

	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, err
	}

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return false, unmarshalError(body, resp.StatusCode)
	}

	defer resp.Body.Close()
	return true, nil
}
