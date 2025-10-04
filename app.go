package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/types"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	summonerEndpoint        = "/lol-summoner/v1/current-summoner"
	profileIconPathTemplate = "/lol-game-data/assets/v1/profile-icons/%d.jpg"
	defaultProfileIconID    = 4804
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	port, token, region, err := detectRiotCredentials()
	if err != nil {
		fmt.Printf("[app.NewApp] failed to detect Riot credentials: %v\n", err)
	} else {
		config.InitInstance(port, token, region)
		fmt.Printf("[app.NewApp] detected Riot client on port %d (region %s)\n", port, region)
	}

	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) types.IPlayerBaseData {
	const logPrefix = "[app.Greet]"

	cfg, ok := config.Instance()
	if !ok {
		fmt.Println(logPrefix + " config not initialised; returning empty payload")
		return types.IPlayerBaseData{}
	}

	resp, err := cfg.SendHttpRequest(summonerEndpoint, http.MethodGet)
	if err != nil {
		fmt.Printf("%s request failed: %v\n", logPrefix, err)
		return types.IPlayerBaseData{}
	}

	fmt.Printf("%s received %d bytes of summoner data\n", logPrefix, len(resp))
	var baseData types.IPlayerBaseData
	if err := json.Unmarshal(resp, &baseData); err != nil {
		fmt.Printf("%s unmarshal failed: %v\n", logPrefix, err)
		return types.IPlayerBaseData{}
	}
	baseData.Region = cfg.Region
	return baseData
}

func (a *App) GetImgSrc(iconID int) string {
	const logPrefix = "[app.GetImgSrc]"

	cfg, ok := config.Instance()
	if !ok {
		fmt.Println(logPrefix + " config not initialised; returning empty result")
		return ""
	}

	icon := iconID
	if icon <= 0 {
		icon = defaultProfileIconID
	}

	path := fmt.Sprintf(profileIconPathTemplate, icon)
	resp, err := cfg.SendHttpRequest(path, http.MethodGet)
	if err != nil {
		fmt.Printf("%s failed to request icon %d: %v\n", logPrefix, icon, err)
		return ""
	}

	fmt.Printf("%s received %d bytes for icon %d\n", logPrefix, len(resp), icon)
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(resp)
}

func detectRiotCredentials() (int, string, string, error) {
	cmd := exec.Command("wmic", "PROCESS", "WHERE", "name='LeagueClientUx.exe'", "GET", "commandline")
	output, err := cmd.Output()
	if err != nil {
		return 0, "", "", err
	}

	text := strings.ReplaceAll(string(output), "\r", "\n")

	port, err := extractPort(text)
	if err != nil {
		return 0, "", "", err
	}

	token, err := extractToken(text)
	if err != nil {
		return 0, "", "", err
	}

	region, err := extractRegion(text)
	if err != nil {
		return 0, "", "", err
	}

	fmt.Printf("[app.detectRiotCredentials] detected port %d (token length %d, region %s)\n", port, len(token), region)
	return port, token, region, nil
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

func extractRegion(text string) (string, error) {
	re := regexp.MustCompile(`(?i)--rso_platform_id[=\s]+([\w-]+)`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return "", errors.New("platform id not found")
	}

	return strings.TrimSpace(match[1]), nil
}
