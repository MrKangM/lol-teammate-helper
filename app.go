package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"lol-teammate-helper/internal/reqLib"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	localHost               = "https://127.0.0.1"
	defaultLocalPort        = 56737
	summonerEndpoint        = "/lol-summoner/v1/current-summoner"
	profileIconPathTemplate = "/lol-game-data/assets/v1/profile-icons/%d.jpg"
	defaultProfileIconID    = 4804
	jsonContentType         = "application/json"
)

type App struct {
	ctx        context.Context
	token      string
	authHeader string
	port       int
}

func NewApp() *App {
	port := defaultLocalPort
	token := ""

	if detectedPort, detectedToken, err := detectRiotCredentials(); err == nil {
		port = detectedPort
		token = detectedToken
	} else {
		fmt.Printf("failed to detect Riot credentials: %v\n", err)
	}

	if token == "" {
		fmt.Println("riot credentials unavailable; requests will likely fail until the League client is running")
	}

	return &App{
		token:      token,
		authHeader: encodeAuthHeader(token),
		port:       port,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	resp, err := a.doGet(summonerEndpoint)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("读取响应失败: %v", err)
	}

	return string(body)
}

func (a *App) GetImgSrc(iconID int) string {
	icon := iconID
	if icon <= 0 {
		icon = defaultProfileIconID
	}

	path := fmt.Sprintf(profileIconPathTemplate, icon)

	resp, err := a.doGet(path)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("请求失败，状态码: %d", resp.StatusCode)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("读取图片数据失败: %v", err)
	}

	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imgData)
}

func (a *App) doGet(path string) (*http.Response, error) {
	req, err := a.newRequest(http.MethodGet, path)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := reqLib.Instance().Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return resp, nil
}

func (a *App) newRequest(method, path string) (*http.Request, error) {
	url := a.buildURL(path)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Accept", jsonContentType)
	req.Header.Set("Authorization", a.authHeader)

	return req, nil
}

func (a *App) buildURL(path string) string {
	return fmt.Sprintf("%s:%d%s", localHost, a.port, path)
}

func encodeAuthHeader(token string) string {
	credential := "riot:" + token
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credential))
}

func (a *App) RiotPort() int {
	return a.port
}

func (a *App) RiotToken() string {
	return a.token
}

func detectRiotCredentials() (int, string, error) {
	cmd := exec.Command("wmic", "PROCESS", "WHERE", "name='LeagueClientUx.exe'", "GET", "commandline")
	output, err := cmd.Output()
	if err != nil {
		return 0, "", err
	}

	text := strings.ReplaceAll(string(output), "\r", "\n")

	port, err := extractPort(text)
	if err != nil {
		return 0, "", err
	}

	token, err := extractToken(text)
	if err != nil {
		return 0, "", err
	}
	fmt.Println("token:", token, "port:", port)
	return port, token, nil
}

func extractPort(text string) (int, error) {
	re := regexp.MustCompile(`(?i)--app-port[=\s]+(\d+)`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return 0, errors.New("app port not found")
	}

	port, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("invalid port value: %w", err)
	}

	return port, nil
}

func extractToken(text string) (string, error) {
	re := regexp.MustCompile(`(?i)--remoting-auth-token[=\s]+([\w-]+)`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return "", errors.New("auth token not found")
	}

	return strings.TrimSpace(match[1]), nil
}
