package dispatch

import (
	"reflect"
	"sync"

	"lol-teammate-helper/internal/types"
)

var champSelectState struct {
	mu       sync.RWMutex
	snapshot types.ChampSelectSnapshot
	ready    bool
}

// StoreChampSelectSnapshot replaces the cached snapshot. It returns false when
// the new snapshot is identical to the stored one (ignoring UpdatedAt), so
// callers can skip redundant UI updates.
func StoreChampSelectSnapshot(snapshot types.ChampSelectSnapshot) bool {
	champSelectState.mu.Lock()
	defer champSelectState.mu.Unlock()

	if champSelectState.ready && sameContent(champSelectState.snapshot, snapshot) {
		return false
	}
	champSelectState.snapshot = snapshot
	champSelectState.ready = true
	return true
}

// ClearChampSelectSnapshot forgets the cached snapshot (champ select finished).
func ClearChampSelectSnapshot() {
	champSelectState.mu.Lock()
	champSelectState.snapshot = types.ChampSelectSnapshot{}
	champSelectState.ready = false
	champSelectState.mu.Unlock()
}

// GetChampSelectSnapshot returns the cached snapshot if available.
func GetChampSelectSnapshot() (types.ChampSelectSnapshot, bool) {
	champSelectState.mu.RLock()
	defer champSelectState.mu.RUnlock()

	if !champSelectState.ready {
		return types.ChampSelectSnapshot{}, false
	}
	return champSelectState.snapshot, true
}

func sameContent(a, b types.ChampSelectSnapshot) bool {
	a.UpdatedAt = b.UpdatedAt
	return reflect.DeepEqual(a, b)
}
