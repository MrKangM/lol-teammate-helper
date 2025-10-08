package service

import (
	"fmt"

	"lol-teammate-helper/internal/types"
)

// TeammateMatchData bundles a teammate's recent match history with champion metadata.
type TeammateMatchData struct {
	History types.MatchHistory
	Heroes  map[int]types.HeroInfo
}

// GetTeammateMatchDetails fetches a teammate's ranked match history plus champion info.
func GetTeammateMatchDetails(puuid string) (TeammateMatchData, error) {
	var result TeammateMatchData

	svc := NewMatchHistoryService()

	matchHistory, err := svc.GetPlayerRankMatches(puuid)
	if err != nil {
		return result, err
	}

	fmt.Printf("[service.GetTeammateMatchDetails] received match history for %s with %d games\n", puuid, len(matchHistory.Games.Games))
	result.History = matchHistory

	championIDs := collectChampionIDs(matchHistory)
	if len(championIDs) == 0 {
		return result, nil
	}

	heroes, err := svc.GetMatchHistoryHeroesByIds(championIDs)
	if err != nil {
		fmt.Printf("[service.GetTeammateMatchDetails] failed to fetch hero info: %v\n", err)
		return result, nil
	}

	result.Heroes = heroes
	return result, nil
}

func collectChampionIDs(history types.MatchHistory) []int {
	games := history.Games.Games
	if len(games) == 0 {
		return nil
	}

	unique := make(map[int]struct{})
	for _, game := range games {
		for _, participant := range game.Participants {
			if participant.ChampionID <= 0 {
				continue
			}
			unique[participant.ChampionID] = struct{}{}
		}
	}

	if len(unique) == 0 {
		return nil
	}

	ids := make([]int, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}

	return ids
}
