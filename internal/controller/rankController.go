package controller

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

type MatchHistory struct {
	cacheMu   sync.RWMutex
	heroCache map[int]HeroInfo
}

type HeroInfo struct {
	Name               string `json:"name"`
	SquarePortraitPath string `json:"squarePortraitPath"`
	IconDataURI        string `json:"iconDataURI,omitempty"`
}

func NewMatchHistory() *MatchHistory {
	return &MatchHistory{
		heroCache: make(map[int]HeroInfo),
	}
}

/*
获取玩家比赛记录数据
*/
func (mh *MatchHistory) GetPlayerRankMatches(puuid string) (types.MatchHistory, error) {
	const logPrefix = "[rankController.GetPlayerRankMatches]"

	if puuid == "" {
		err := errors.New("empty puuid provided")
		fmt.Printf("%s %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	cfg, ok := config.Instance()
	if !ok {
		err := errors.New("riot credentials are not initialised")
		fmt.Printf("%s %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	endpoint := fmt.Sprintf("/lol-match-history/v1/products/lol/%s/matches", puuid)
	fmt.Println(logPrefix + " requesting " + endpoint)

	body, err := cfg.SendHttpRequest(endpoint, http.MethodGet)
	if err != nil {
		fmt.Printf("%s failed to request matches: %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	fmt.Printf("%s received %d bytes\n", logPrefix, len(body))

	var matchRecord types.MatchHistory
	if err = json.Unmarshal(body, &matchRecord); err != nil {
		fmt.Printf("%s failed to unmarshal response: %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	return matchRecord, nil
}

/*
获取英雄数据，并在内存中缓存减轻重复请求压力
*/
func (mh *MatchHistory) GetMatchHistoryNameAndIconByHeroId(id int) (HeroInfo, error) {
	if id <= 0 {
		return HeroInfo{}, errors.New("invalid id provided")
	}

	heroes, err := mh.GetMatchHistoryHeroesByIds([]int{id})
	if err != nil {
		return HeroInfo{}, err
	}

	info, ok := heroes[id]
	if !ok {
		return HeroInfo{}, fmt.Errorf("hero %d not found", id)
	}

	return info, nil
}

func (mh *MatchHistory) GetMatchHistoryHeroesByIds(ids []int) (map[int]HeroInfo, error) {
	result := make(map[int]HeroInfo)
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

		if info, ok := mh.getHeroFromCache(id); ok {
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
		info, err := mh.fetchHeroInfo(cfg, id)
		if err != nil {
			return nil, err
		}
		mh.storeHeroInCache(id, info)
		result[id] = info
	}

	return result, nil
}

func (mh *MatchHistory) fetchHeroInfo(cfg *config.AppConfig, id int) (HeroInfo, error) {
	url := fmt.Sprintf(heroByIDPath, id)

	resp, err := cfg.SendHttpRequest(url, http.MethodGet)
	if err != nil {
		return HeroInfo{}, err
	}

	var heroInfo HeroInfo
	if err = json.Unmarshal(resp, &heroInfo); err != nil {
		return HeroInfo{}, err
	}

	baseURL := fmt.Sprintf(baseURLTemplate, cfg.Port)
	heroIconURL := fmt.Sprintf("%s%s", baseURL, heroInfo.SquarePortraitPath)
	heroInfo.SquarePortraitPath = heroIconURL

	if iconBytes, iconErr := cfg.SendHttpRequest(heroInfo.SquarePortraitPath, http.MethodGet); iconErr == nil && len(iconBytes) > 0 {
		heroInfo.IconDataURI = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(iconBytes))
	} else if iconErr != nil {
		fmt.Printf("[rankController.fetchHeroInfo] failed to download icon %s: %v\n", heroInfo.SquarePortraitPath, iconErr)
	}

	return heroInfo, nil
}

func (mh *MatchHistory) getHeroFromCache(id int) (HeroInfo, bool) {
	if mh == nil {
		return HeroInfo{}, false
	}

	mh.cacheMu.RLock()
	defer mh.cacheMu.RUnlock()

	info, ok := mh.heroCache[id]
	return info, ok
}

func (mh *MatchHistory) storeHeroInCache(id int, info HeroInfo) {
	if mh == nil {
		return
	}

	mh.cacheMu.Lock()
	defer mh.cacheMu.Unlock()

	if mh.heroCache == nil {
		mh.heroCache = make(map[int]HeroInfo)
	}
	mh.heroCache[id] = info
}
