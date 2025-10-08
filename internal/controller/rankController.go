package controller

import (
	"errors"
	"fmt"

	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
)

type MatchHistory struct {
	svc *service.MatchHistoryService
}

func NewMatchHistory() *MatchHistory {
	return &MatchHistory{svc: service.NewMatchHistoryService()}
}

// GetPlayerRankMatches validates the input and delegates to the service layer.
func (mh *MatchHistory) GetPlayerRankMatches(puuid string) (types.MatchHistory, error) {
	const logPrefix = "[rankController.GetPlayerRankMatches]"

	if puuid == "" {
		err := errors.New("empty puuid provided")
		fmt.Printf("%s %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	if mh == nil || mh.svc == nil {
		err := errors.New("match history service not initialised")
		fmt.Printf("%s %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	matches, err := mh.svc.GetPlayerRankMatches(puuid)
	if err != nil {
		fmt.Printf("%s failed to fetch matches: %v\n", logPrefix, err)
		return types.MatchHistory{}, err
	}

	return matches, nil
}

// GetMatchHistoryHeroesByIds validates IDs and defers to the service.
func (mh *MatchHistory) GetMatchHistoryHeroesByIds(ids []int) (map[int]types.HeroInfo, error) {
	const logPrefix = "[rankController.GetMatchHistoryHeroesByIds]"

	if len(ids) == 0 {
		return map[int]types.HeroInfo{}, nil
	}

	if mh == nil || mh.svc == nil {
		return nil, errors.New("match history service not initialised")
	}

	heroes, err := mh.svc.GetMatchHistoryHeroesByIds(ids)
	if err != nil {
		fmt.Printf("%s failed: %v\n", logPrefix, err)
		return nil, err
	}

	return heroes, nil
}

// GetMatchHistoryNameAndIconByHeroId validates the ID and forwards to the service.
func (mh *MatchHistory) GetMatchHistoryNameAndIconByHeroId(id int) (types.HeroInfo, error) {
	const logPrefix = "[rankController.GetMatchHistoryNameAndIconByHeroId]"

	if id <= 0 {
		err := errors.New("invalid id provided")
		fmt.Printf("%s %v\n", logPrefix, err)
		return types.HeroInfo{}, err
	}

	if mh == nil || mh.svc == nil {
		return types.HeroInfo{}, errors.New("match history service not initialised")
	}

	info, err := mh.svc.GetMatchHistoryNameAndIconByHeroId(id)
	if err != nil {
		fmt.Printf("%s failed: %v\n", logPrefix, err)
		return types.HeroInfo{}, err
	}

	return info, nil
}
