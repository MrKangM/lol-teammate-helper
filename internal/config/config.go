package config

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"lol-teammate-helper/internal/utils"
)

const BaseURL = "https://127.0.0.1"

// AppConfig holds the connection details of the running League client.
// Instances returned by Instance are immutable snapshots; use Update to change them.
type AppConfig struct {
	Port      int    `json:"port"`
	Token     string `json:"token"`
	MetaToken string `json:"meta_token"`
	Region    string `json:"region"`
}

var (
	mu       sync.RWMutex
	instance *AppConfig

	// The LCU only listens on loopback with a self-signed certificate.
	riotTransport = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        32,
		MaxIdleConnsPerHost: 8,
		IdleConnTimeout:     90 * time.Second,
	}
	riotHTTPClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: riotTransport,
	}
)

// Update replaces the stored credentials. It returns true when anything changed
// (the League client picks a new port and token on every launch).
func Update(port int, rawToken string, region string) bool {
	next := &AppConfig{
		Port:      port,
		Token:     "Basic " + base64.StdEncoding.EncodeToString([]byte("riot:"+rawToken)),
		MetaToken: rawToken,
		Region:    utils.GetServerChineseName(region),
	}

	mu.Lock()
	defer mu.Unlock()
	if instance != nil && *instance == *next {
		return false
	}
	instance = next
	slog.Info("config updated", "port", port, "region", next.Region)
	return true
}

// Instance returns the current credentials, if the client has been detected.
func Instance() (*AppConfig, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if instance == nil {
		return nil, false
	}
	snapshot := *instance
	return &snapshot, true
}

// SendHttpRequest issues a request against the LCU REST API (or an absolute URL).
func (ac *AppConfig) SendHttpRequest(endpoint string, method string) ([]byte, error) {
	if ac == nil {
		return nil, fmt.Errorf("app config is nil")
	}
	if ac.Port <= 0 {
		return nil, fmt.Errorf("invalid port: %d", ac.Port)
	}
	if ac.Token == "" {
		return nil, fmt.Errorf("authorization token is empty")
	}

	url := endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		url = fmt.Sprintf("%s:%d%s", BaseURL, ac.Port, endpoint)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	resp, err := riotHTTPClient.Do(req)
	if err != nil {
		slog.Warn("lcu request failed", "method", method, "endpoint", endpoint, "err", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		slog.Warn("lcu request rejected", "method", method, "endpoint", endpoint, "status", resp.Status)
		return nil, fmt.Errorf("http request %s failed with status: %s", endpoint, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	slog.Debug("lcu request", "method", method, "endpoint", endpoint, "bytes", len(body))
	return body, nil
}
