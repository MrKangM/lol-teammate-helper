package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/types"
)

const (
	baseURLTemplate = "https://127.0.0.1:%d"
	heroByIDPath    = "/lol-game-data/assets/v1/champions/%d.json"
	matchesPath     = "/lol-match-history/v1/products/lol/%s/matches"

	// matchCacheTTL bounds how stale a teammate's history may be. Champ select
	// emits many updates per lobby, so this avoids re-querying every time.
	matchCacheTTL = 3 * time.Minute
)

type matchCacheEntry struct {
	history   types.MatchHistory
	expiresAt time.Time
}

// MatchHistoryService encapsulates the business logic and caching layer for match history calls.
type MatchHistoryService struct {
	cacheMu    sync.RWMutex
	heroCache  map[int]types.HeroInfo
	matchCache map[string]matchCacheEntry
	now        func() time.Time

	// Generic LCU response bodies (rank, level, mastery) keyed by endpoint.
	bodyCache map[string]bodyCacheEntry
	// Ranked game lists keyed by puuid.
	gamesCache map[string]gamesCacheEntry
	// Item metadata for the end-of-game screen.
	itemPaths map[int]string
	itemIcons map[int]string
}

// NewMatchHistoryService constructs a service with empty caches.
func NewMatchHistoryService() *MatchHistoryService {
	return &MatchHistoryService{
		heroCache:  make(map[int]types.HeroInfo),
		matchCache: make(map[string]matchCacheEntry),
		bodyCache:  make(map[string]bodyCacheEntry),
		gamesCache: make(map[string]gamesCacheEntry),
		itemIcons:  make(map[int]string),
		now:        time.Now,
	}
}

var (
	sharedOnce sync.Once
	shared     *MatchHistoryService
)

// Shared returns the process-wide service so every caller benefits from the same caches.
func Shared() *MatchHistoryService {
	sharedOnce.Do(func() { shared = NewMatchHistoryService() })
	return shared
}

// ResetCaches drops cached data; call it when the client reconnects (new session, new assets).
func (svc *MatchHistoryService) ResetCaches() {
	svc.cacheMu.Lock()
	defer svc.cacheMu.Unlock()
	svc.heroCache = make(map[int]types.HeroInfo)
	svc.matchCache = make(map[string]matchCacheEntry)
	svc.bodyCache = make(map[string]bodyCacheEntry)
	svc.gamesCache = make(map[string]gamesCacheEntry)
	svc.itemPaths = nil
	svc.itemIcons = make(map[int]string)
}

// GetPlayerRankMatches returns the ranked match history for the provided PUUID.
func (svc *MatchHistoryService) GetPlayerRankMatches(puuid string) (types.MatchHistory, error) {
	if h, ok := svc.getMatchesFromCache(puuid); ok {
		return h, nil
	}

	cfg, ok := config.Instance()
	if !ok {
		return types.MatchHistory{}, errors.New("riot credentials are not initialised")
	}

	body, err := cfg.SendHttpRequest(fmt.Sprintf(matchesPath, puuid), http.MethodGet)
	if err != nil {
		return types.MatchHistory{}, err
	}

	var matchRecord types.MatchHistory
	if err = json.Unmarshal(body, &matchRecord); err != nil {
		return types.MatchHistory{}, err
	}

	svc.cacheMu.Lock()
	svc.matchCache[puuid] = matchCacheEntry{history: matchRecord, expiresAt: svc.now().Add(matchCacheTTL)}
	svc.cacheMu.Unlock()
	return matchRecord, nil
}

// GetMatchHistoryHeroesByIds fetches hero metadata, using the internal cache when possible.
// Heroes that fail to load are skipped so one bad lookup does not discard the rest;
// an error is only returned when nothing could be resolved.
func (svc *MatchHistoryService) GetMatchHistoryHeroesByIds(ids []int) (map[int]types.HeroInfo, error) {
	result := make(map[int]types.HeroInfo)
	seen := make(map[int]struct{}, len(ids))
	var missing []int

	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}

		if info, ok := svc.getHeroFromCache(id); ok {
			result[id] = info
		} else {
			missing = append(missing, id)
		}
	}
	if len(missing) == 0 {
		return result, nil
	}

	cfg, ok := config.Instance()
	if !ok {
		return nil, errors.New("riot credentials are not initialised")
	}

	var (
		wg       sync.WaitGroup
		resMu    sync.Mutex
		firstErr error
	)
	for _, id := range missing {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			info, err := svc.fetchHeroInfo(cfg, id)
			resMu.Lock()
			defer resMu.Unlock()
			if err != nil {
				slog.Warn("fetch hero failed", "id", id, "err", err)
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			svc.storeHeroInCache(id, info)
			result[id] = info
		}(id)
	}
	wg.Wait()

	if len(result) == 0 && firstErr != nil {
		return nil, firstErr
	}
	return result, nil
}

// GetMatchHistoryNameAndIconByHeroId returns metadata for a single hero.
func (svc *MatchHistoryService) GetMatchHistoryNameAndIconByHeroId(id int) (types.HeroInfo, error) {
	heroes, err := svc.GetMatchHistoryHeroesByIds([]int{id})
	if err != nil {
		return types.HeroInfo{}, err
	}
	info, ok := heroes[id]
	if !ok {
		return types.HeroInfo{}, fmt.Errorf("hero %d not found", id)
	}
	return info, nil
}

func (svc *MatchHistoryService) fetchHeroInfo(cfg *config.AppConfig, id int) (types.HeroInfo, error) {
	resp, err := cfg.SendHttpRequest(fmt.Sprintf(heroByIDPath, id), http.MethodGet)
	if err != nil {
		return types.HeroInfo{}, err
	}

	var heroInfo types.HeroInfo
	if err = json.Unmarshal(resp, &heroInfo); err != nil {
		return types.HeroInfo{}, err
	}

	heroInfo.SquarePortraitPath = fmt.Sprintf(baseURLTemplate, cfg.Port) + heroInfo.SquarePortraitPath

	iconBytes, iconErr := cfg.SendHttpRequest(heroInfo.SquarePortraitPath, http.MethodGet)
	if iconErr != nil {
		slog.Warn("download hero icon failed", "path", heroInfo.SquarePortraitPath, "err", iconErr)
	} else if len(iconBytes) > 0 {
		heroInfo.IconDataURI = "data:image/png;base64," + base64.StdEncoding.EncodeToString(iconBytes)
	}
	return heroInfo, nil
}

func (svc *MatchHistoryService) getMatchesFromCache(puuid string) (types.MatchHistory, bool) {
	svc.cacheMu.RLock()
	defer svc.cacheMu.RUnlock()
	entry, ok := svc.matchCache[puuid]
	if !ok || !svc.now().Before(entry.expiresAt) {
		return types.MatchHistory{}, false
	}
	return entry.history, true
}

func (svc *MatchHistoryService) getHeroFromCache(id int) (types.HeroInfo, bool) {
	svc.cacheMu.RLock()
	defer svc.cacheMu.RUnlock()
	info, ok := svc.heroCache[id]
	return info, ok
}

func (svc *MatchHistoryService) storeHeroInCache(id int, info types.HeroInfo) {
	svc.cacheMu.Lock()
	defer svc.cacheMu.Unlock()
	svc.heroCache[id] = info
}
