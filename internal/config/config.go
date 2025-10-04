package config

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"lol-teammate-helper/internal/utils"
	"net/http"
	"strings"
	"sync"
	"time"
)

const BaseURL = "https://127.0.0.1"

const (
	initLogPrefix    = "[config.InitInstance]"
	requestLogPrefix = "[config.SendHttpRequest]"
)

type AppConfig struct {
	Port      int    `json:"port"`
	Token     string `json:"token"`
	MetaToken string `json:"meta_token"`
	Region    string `json:"region"`
}

var (
	instance *AppConfig
	once     sync.Once
)

func InitInstance(port int, token string, region string) {
	once.Do(func() {
		authString := "riot:" + token
		chineseRegion := utils.GetServerChineseName(region)
		instance = &AppConfig{
			Port:      port,
			Token:     "Basic " + base64.StdEncoding.EncodeToString([]byte(authString)),
			MetaToken: token,
			Region:    chineseRegion,
		}
		fmt.Printf("%s initialised config with port %d (region %s)\n", initLogPrefix, port, chineseRegion)
	})
}

func GetInstance() *AppConfig {
	if instance == nil {
		panic("config instance not initialised")
	}
	return instance
}

func Instance() (*AppConfig, bool) {
	if instance == nil {
		return nil, false
	}
	return instance, true
}

func (ac *AppConfig) SendHttpRequest(endpoint string, method string) ([]byte, error) {
	if ac == nil {
		fmt.Println(requestLogPrefix + " app config is nil")
		return nil, fmt.Errorf("app config is nil")
	}

	if ac.Port <= 0 {
		fmt.Printf("%s invalid port: %d\n", requestLogPrefix, ac.Port)
		return nil, fmt.Errorf("invalid port: %d", ac.Port)
	}
	if ac.Token == "" {
		fmt.Println(requestLogPrefix + " authorization token is empty")
		return nil, fmt.Errorf("authorization token is empty")
	}

	url := endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		url = fmt.Sprintf("%s:%d%s", BaseURL, ac.Port, endpoint)
	}

	fmt.Printf("%s sending %s %s\n", requestLogPrefix, method, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Printf("%s failed to create request: %v\n", requestLogPrefix, err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s failed to send request: %v\n", requestLogPrefix, err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("%s received status %d\n", requestLogPrefix, resp.StatusCode)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		fmt.Printf("%s http request failed with status: %d %s\n", requestLogPrefix, resp.StatusCode, resp.Status)
		return nil, fmt.Errorf("http request failed with status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s failed to read response body: %v\n", requestLogPrefix, err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("%s received %d bytes\n", requestLogPrefix, len(body))

	return body, nil
}
