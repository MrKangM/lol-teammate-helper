package dispatch

import (
	"context"
	"encoding/json"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
)

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

// LCU resources routed by EventHandler.
const (
	ChampSelectURI = "/lol-champ-select/v1/session"
	GameflowURI    = "/lol-gameflow/v1/session"
)

// EventHandler routes websocket events based on their URL. It must be called
// from a single goroutine: handlers read-modify-write the shared snapshot.
func EventHandler(uri string, eventType string, data json.RawMessage) {
	switch uri {
	case ChampSelectURI:
		if eventType == "Delete" {
			return // the gameflow phase decides when the snapshot is cleared
		}
		handleChampSelectEvent(data)
	case GameflowURI:
		if eventType == "Delete" {
			handleSessionEnded()
			return
		}
		handleGameflowEvent(data)
	default:
		slog.Debug("unhandled event", "uri", uri)
	}
}

func emit(name string, payload ...interface{}) {
	if ctx := getRuntimeContext(); ctx != nil {
		runtime.EventsEmit(ctx, name, payload...)
	}
}

func handleSessionEnded() {
	ClearChampSelectSnapshot()
	emit(EventEnded)
}

// publish stores the snapshot and notifies the UI when something changed.
func publish(snapshot types.ChampSelectSnapshot) {
	snapshot.UpdatedAt = time.Now()
	if StoreChampSelectSnapshot(snapshot) {
		emit(EventSnapshot, snapshot)
	}
}

func handleChampSelectEvent(data json.RawMessage) {
	var cs types.ChampSelectData
	if err := json.Unmarshal(data, &cs); err != nil {
		slog.Warn("decode champ select payload failed", "err", err)
		return
	}

	own := make([]memberInput, len(cs.MyTeam))
	for i, p := range cs.MyTeam {
		own[i] = inputFromPlayer(p)
	}
	team := buildMembers(own)
	sortByCell(team)

	// Ranked champ select hides the opposing puuids; custom/normal modes may not.
	var enemyInputs []memberInput
	for _, p := range cs.TheirTeam {
		if p.Puuid != "" {
			enemyInputs = append(enemyInputs, inputFromPlayer(p))
		}
	}

	snapshot := types.ChampSelectSnapshot{
		QueueID: cs.QueueID,
		GameID:  cs.GameID,
		Phase:   "ChampSelect",
		Team:    team,
	}
	if len(enemyInputs) > 0 {
		snapshot.Enemy = buildMembers(enemyInputs)
		sortByCell(snapshot.Enemy)
	}
	publish(snapshot)
}

func sortByCell(members []types.TeamMemberSummary) {
	sort.SliceStable(members, func(i, j int) bool { return members[i].CellID < members[j].CellID })
}
