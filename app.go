package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/connector"
	"lol-teammate-helper/internal/dispatch"
	"lol-teammate-helper/internal/types"
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
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dispatch.SetRuntimeContext(ctx)
	go connector.Run(ctx)
}

// GetCurrentSummoner returns the logged-in summoner's profile, or an empty
// payload while the League client has not been detected yet.
func (a *App) GetCurrentSummoner() types.IPlayerBaseData {
	cfg, ok := config.Instance()
	if !ok {
		return types.IPlayerBaseData{}
	}

	resp, err := cfg.SendHttpRequest(summonerEndpoint, http.MethodGet)
	if err != nil {
		slog.Warn("fetch current summoner failed", "err", err)
		return types.IPlayerBaseData{}
	}

	var baseData types.IPlayerBaseData
	if err := json.Unmarshal(resp, &baseData); err != nil {
		slog.Warn("decode current summoner failed", "err", err)
		return types.IPlayerBaseData{}
	}
	baseData.Region = cfg.Region
	return baseData
}

// GetImgSrc returns a profile icon as a data URI.
func (a *App) GetImgSrc(iconID int) string {
	cfg, ok := config.Instance()
	if !ok {
		return ""
	}

	if iconID <= 0 {
		iconID = defaultProfileIconID
	}

	resp, err := cfg.SendHttpRequest(profileIconURL(iconID), http.MethodGet)
	if err != nil {
		slog.Warn("fetch profile icon failed", "icon", iconID, "err", err)
		return ""
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(resp)
}

// GetCurrentChampSelectSnapshot exposes the latest cached champion select snapshot to the frontend.
func (a *App) GetCurrentChampSelectSnapshot() types.ChampSelectSnapshot {
	snapshot, ok := dispatch.GetChampSelectSnapshot()
	if !ok {
		return types.ChampSelectSnapshot{}
	}
	return snapshot
}

func profileIconURL(id int) string {
	return fmt.Sprintf(profileIconPathTemplate, id)
}
