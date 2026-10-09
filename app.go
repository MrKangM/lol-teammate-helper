package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"runtime"

	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/connector"
	"lol-teammate-helper/internal/diag"
	"lol-teammate-helper/internal/dispatch"
	"lol-teammate-helper/internal/logging"
	"lol-teammate-helper/internal/service"
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

// GetMyCareer analyses the logged-in summoner's own recent ranked games.
func (a *App) GetMyCareer() types.TeamMemberSummary {
	me, ok := dispatch.BuildSelf()
	if !ok {
		return types.TeamMemberSummary{}
	}
	return me
}

// GetRankEmblem returns the client's own emblem image for a tier ("" if unavailable).
func (a *App) GetRankEmblem(tierKey string) string {
	return service.Shared().GetAsset(service.RankEmblemCandidates(tierKey))
}

// GetPositionIcon returns the client's own lane icon ("" if unavailable).
func (a *App) GetPositionIcon(position string) string {
	if position == "" {
		return ""
	}
	return service.Shared().GetAsset(service.PositionIconCandidates(position))
}

// GetDiagnostics returns connection state and the most recent raw events.
func (a *App) GetDiagnostics() diag.Snapshot {
	return diag.Get()
}

// OpenLogDir opens the folder containing the log file in the file manager.
func (a *App) OpenLogDir() {
	dir := logging.Dir()
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	// explorer reports a non-zero exit code even on success, so only start it.
	if err := cmd.Start(); err != nil {
		slog.Warn("open log dir failed", "dir", dir, "err", err)
	}
}
