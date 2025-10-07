package service

import (
	"fmt"

	"lol-teammate-helper/internal/controller"
	"lol-teammate-helper/internal/types"
)

// TeammateMatchData bundles a teammate\'s recent match history with the hero metadata
// required by the UI.
type TeammateMatchData struct {
	History types.MatchHistory
	Heroes  map[int]controller.HeroInfo
}

// GetTeammateMatchDetails fetches the teammate\'s ranked match history along with the
// hero metadata needed to render champion icons.
func GetTeammateMatchDetails(puuid string) (TeammateMatchData, error) {
	var result TeammateMatchData

	mhController := controller.NewMatchHistory()
	matchHistory, err := mhController.GetPlayerRankMatches(puuid)
	if err != nil {
		return result, err
	}

	fmt.Printf("GetTeammateMatchHistory杩斿洖鐨勭帺瀹舵暟鎹細%+v\n", matchHistory)
	fmt.Println("GetTeammateMatchHistory鑾峰彇鐨刾uuid" + puuid)

	result.History = matchHistory

	championIDs := collectChampionIDs(matchHistory)
	if len(championIDs) == 0 {
		return result, nil
	}

	heroes, err := mhController.GetMatchHistoryHeroesByIds(championIDs)
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
