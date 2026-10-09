package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
)

const maxRecentMatches = 5

var (
	runtimeCtx      context.Context
	runtimeCtxMu    sync.RWMutex
	matchHistorySvc = service.Shared()
)

// SetRuntimeContext allows the Wails runtime context to be reused when emitting events.
func SetRuntimeContext(ctx context.Context) {
	runtimeCtxMu.Lock()
	defer runtimeCtxMu.Unlock()
	runtimeCtx = ctx
}

func getRuntimeContext() context.Context {
	runtimeCtxMu.RLock()
	defer runtimeCtxMu.RUnlock()
	return runtimeCtx
}

// Event names emitted to the frontend.
const (
	EventSnapshot = "champ-select:snapshot"
	EventEnded    = "champ-select:ended"
)

const champSelectURI = "/lol-champ-select/v1/session"

// EventHandler routes websocket events based on their URL.
func EventHandler(uri string, eventType string, data json.RawMessage) {
	switch uri {
	case champSelectURI:
		if eventType == "Delete" {
			handleChampSelectEnded()
			return
		}
		handleChampSelectEvent(data)
	default:
		slog.Debug("unhandled event", "uri", uri)
	}
}

func emit(name string, payload ...interface{}) {
	if ctx := getRuntimeContext(); ctx != nil {
		runtime.EventsEmit(ctx, name, payload...)
	}
}

func handleChampSelectEnded() {
	ClearChampSelectSnapshot()
	emit(EventEnded)
}

func handleChampSelectEvent(data json.RawMessage) {
	var champSelect types.ChampSelectData
	if err := json.Unmarshal(data, &champSelect); err != nil {
		slog.Warn("decode champ select payload failed", "err", err)
		return
	}

	snapshot := buildSnapshot(champSelect)
	if !StoreChampSelectSnapshot(snapshot) {
		return // nothing visible changed since the last emit
	}
	emit(EventSnapshot, snapshot)
}

func buildSnapshot(champSelect types.ChampSelectData) types.ChampSelectSnapshot {
	summaries := make([]types.TeamMemberSummary, len(champSelect.MyTeam))

	var wg sync.WaitGroup
	for i, player := range champSelect.MyTeam {
		wg.Add(1)
		go func(i int, p types.Player) {
			defer wg.Done()
			summaries[i] = buildTeamMemberSummary(p)
		}(i, player)
	}
	wg.Wait()

	sort.SliceStable(summaries, func(i, j int) bool {
		return summaries[i].CellID < summaries[j].CellID
	})

	return types.ChampSelectSnapshot{
		QueueID:   champSelect.QueueID,
		GameID:    champSelect.GameID,
		UpdatedAt: time.Now(),
		Team:      summaries,
	}
}

func buildTeamMemberSummary(player types.Player) types.TeamMemberSummary {
	summary := types.TeamMemberSummary{
		Puuid:            player.Puuid,
		GameName:         player.GameName,
		TagLine:          player.TagLine,
		AssignedPosition: player.AssignedPosition,
		ChampionID:       player.ChampionID,
		CellID:           player.CellID,
	}

	if player.Puuid == "" {
		return summary
	}

	matchData, err := service.GetTeammateMatchDetails(player.Puuid)
	if err != nil {
		slog.Warn("fetch match data failed", "puuid", player.Puuid, "err", err)
		return summary
	}

	summary.RecentMatches = buildRecentMatches(matchData.History, matchData.Heroes)
	if player.ChampionID > 0 && matchHistorySvc != nil {
		heroInfo, heroErr := matchHistorySvc.GetMatchHistoryNameAndIconByHeroId(player.ChampionID)
		if heroErr != nil {
			slog.Warn("fetch current hero failed", "id", player.ChampionID, "err", heroErr)
		} else {
			summary.ChampionName = heroInfo.Name
			icon := heroInfo.IconDataURI
			if icon == "" {
				icon = heroInfo.SquarePortraitPath
			}
			summary.ChampionIcon = icon
		}
	}
	if summary.ChampionName == "" && player.ChampionID > 0 {
		summary.ChampionName = fmt.Sprintf("Champion %d", player.ChampionID)
	}

	return summary
}

// buildRecentMatches reads the queried player's own stats, which the LCU match
// history endpoint always returns as the first (and only) participant of a game.
func buildRecentMatches(history types.MatchHistory, heroMap map[int]types.HeroInfo) []types.RecentMatchSummary {
	games := history.Games.Games
	if len(games) == 0 {
		return nil
	}

	result := make([]types.RecentMatchSummary, 0, maxRecentMatches)
	for _, game := range games {
		if len(game.Participants) == 0 {
			continue
		}

		participant := game.Participants[0]
		stats := participant.Stats
		heroInfo := heroMap[participant.ChampionID]

		icon := heroInfo.IconDataURI
		if icon == "" {
			icon = heroInfo.SquarePortraitPath
		}

		matchSummary := types.RecentMatchSummary{
			ChampionID:   participant.ChampionID,
			ChampionName: heroInfo.Name,
			ChampionIcon: icon,
			Win:          stats.Win,
			Kills:        stats.Kills,
			Deaths:       stats.Deaths,
			Assists:      stats.Assists,
			QueueID:      game.QueueID,
			GameDuration: game.GameDuration,
		}

		result = append(result, matchSummary)
		if len(result) >= maxRecentMatches {
			break
		}
	}

	return result
}
