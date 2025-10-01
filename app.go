package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"lol-teammate-helper/internal/reqLib"
	"net/http"
)

const (
	localHost               = "https://127.0.0.1"
	localPort               = 56737
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
	token := "iOkyLnA-zyiJNOB9cQ1svA"

	return &App{
		token:      token,
		authHeader: encodeAuthHeader(token),
		port:       localPort,
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
