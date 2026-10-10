package controller

import (
	"errors"
	"log/slog"

	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
)

// MatchHistory exposes match lookups to the frontend.
type MatchHistory struct {
	svc *service.MatchHistoryService
}

func NewMatchHistory() *MatchHistory {
	return &MatchHistory{svc: service.Shared()}
}

// GetGameDetail returns the end-of-game scoreboard (both teams) of one match.
// highlightPuuid marks the row of the player the user is inspecting.
func (mh *MatchHistory) GetGameDetail(gameID int64, highlightPuuid string) (types.GameDetail, error) {
	if gameID <= 0 {
		return types.GameDetail{}, errors.New("invalid game id")
	}
	if mh == nil || mh.svc == nil {
		return types.GameDetail{}, errors.New("match history service not initialised")
	}

	detail, err := mh.svc.GetGameDetail(gameID, highlightPuuid)
	if err != nil {
		slog.Warn("load game detail failed", "game", gameID, "err", err)
		return types.GameDetail{}, err
	}
	return detail, nil
}
