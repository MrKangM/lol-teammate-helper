package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/types"
	"net/http"
	"sync"
)

const (
	baseURLTemplate = "https://127.0.0.1:%d"
	heroByIDPath    = "/lol-game-data/assets/v1/champions/%d.json"
)

// MatchHistoryService encapsulates the business logic and caching layer for match history calls.
type MatchHistoryService struct {
	cacheMu   sync.RWMutex
	heroCache map[int]types.HeroInfo
}

// NewMatchHistoryService constructs a service with an empty hero cache.
func NewMatchHistoryService() *MatchHistoryService {
	return &MatchHistoryService{
		heroCache: make(map[int]types.HeroInfo),
	}
}

// GetPlayerRankMatches returns the ranked match history for the provided PUUID.
func (svc *MatchHistoryService) GetPlayerRankMatches(puuid string) (types.MatchHistory, error) {
	cfg, ok := config.Instance()
	if !ok {
		return types.MatchHistory{}, errors.New("riot credentials are not initialised")
	}

	endpoint := fmt.Sprintf("/lol-match-history/v1/products/lol/%s/matches", puuid)
	fmt.Printf("[service.MatchHistory] requesting %s\n", endpoint)

	body, err := cfg.SendHttpRequest(endpoint, http.MethodGet)
	if err != nil {
		return types.MatchHistory{}, err
	}

	fmt.Printf("[service.MatchHistory] received %d bytes for %s\n", len(body), puuid)

	var matchRecord types.MatchHistory
	if err = json.Unmarshal(body, &matchRecord); err != nil {
		return types.MatchHistory{}, err
	}

	return matchRecord, nil
}

// GetMatchHistoryHeroesByIds fetches hero metadata, using the internal cache when possible.
func (svc *MatchHistoryService) GetMatchHistoryHeroesByIds(ids []int) (map[int]types.HeroInfo, error) {
	result := make(map[int]types.HeroInfo)
	if len(ids) == 0 {
		return result, nil
	}

	unique := make(map[int]struct{}, len(ids))
	missing := make([]int, 0, len(ids))

	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, seen := unique[id]; seen {
			continue
		}
		unique[id] = struct{}{}

		if info, ok := svc.getHeroFromCache(id); ok {
			result[id] = info
			continue
		}
		missing = append(missing, id)
	}

	if len(missing) == 0 {
		return result, nil
	}

	cfg, ok := config.Instance()
	if !ok {
		return nil, errors.New("riot credentials are not initialised")
	}

	for _, id := range missing {
		info, err := svc.fetchHeroInfo(cfg, id)
		if err != nil {
			return nil, err
		}
		svc.storeHeroInCache(id, info)
		result[id] = info
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
	url := fmt.Sprintf(heroByIDPath, id)

	resp, err := cfg.SendHttpRequest(url, http.MethodGet)
	if err != nil {
		return types.HeroInfo{}, err
	}

	var heroInfo types.HeroInfo
	if err = json.Unmarshal(resp, &heroInfo); err != nil {
		return types.HeroInfo{}, err
	}

	baseURL := fmt.Sprintf(baseURLTemplate, cfg.Port)
	heroIconURL := fmt.Sprintf("%s%s", baseURL, heroInfo.SquarePortraitPath)
	heroInfo.SquarePortraitPath = heroIconURL

	if iconBytes, iconErr := cfg.SendHttpRequest(heroInfo.SquarePortraitPath, http.MethodGet); iconErr == nil && len(iconBytes) > 0 {
		heroInfo.IconDataURI = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(iconBytes))
	} else if iconErr != nil {
		fmt.Printf("[service.MatchHistory] failed to download icon %s: %v\n", heroInfo.SquarePortraitPath, iconErr)
	}

	return heroInfo, nil
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

	if svc.heroCache == nil {
		svc.heroCache = make(map[int]types.HeroInfo)
	}
	svc.heroCache[id] = info
}
