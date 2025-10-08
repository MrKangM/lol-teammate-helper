package dispatch

import (
	"sync"

	"lol-teammate-helper/internal/types"
)

var champSelectState struct {
	mu       sync.RWMutex
	snapshot types.ChampSelectSnapshot
	ready    bool
}

// StoreChampSelectSnapshot replaces the cached champion select snapshot.
func StoreChampSelectSnapshot(snapshot types.ChampSelectSnapshot) {
	champSelectState.mu.Lock()
	champSelectState.snapshot = snapshot
	champSelectState.ready = true
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
