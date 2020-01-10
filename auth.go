package dfuse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var AuthEndpoint = "https://auth.dfuse.io"

type TokenStore interface {
	GetAuth() *Authorization
	SetAuth(*Authorization)
}

type InMemoryTokenStore struct {
	auth *Authorization
}

func (ts *InMemoryTokenStore) GetAuth() *Authorization {
	return ts.auth
}

func (ts *InMemoryTokenStore) SetAuth(auth *Authorization) {
	ts.auth = auth
}

func fetchAuth(apiKey string) (*Authorization, error) {
	if apiKey == "" {
		return nil, errors.New("apikey is null")
	}

	u := fmt.Sprintf("%s/v1/auth/issue", AuthEndpoint)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader([]byte(fmt.Sprintf(`{"api_key":"%s"}`, apiKey))))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code = %d", resp.StatusCode)
	}

	var result Authorization
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type Authorization struct {
	Token     string    `json:"token"`
	ExpiresAt Timestamp `json:"expires_at"`
}

func (a *Authorization) IsExpired() bool {
	return a.ExpiresAt.After(time.Now())
}

func SetAuthEndpoint(endpoint string) {
	AuthEndpoint = endpoint
}
