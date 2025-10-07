package dispatch

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"lol-teammate-helper/internal/controller"
	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
	"sync"
)

const maxRecentMatches = 5

// EventHandler routes websocket events based on their URL.
func EventHandler(url string, data json.RawMessage) {
	fmt.Printf("Received event: %s\n", url)
	switch url {
	case "/lol-lobby/v2/lobby":
		fmt.Println("Lobby event")
	case "/lol-champ-select/v1/session":
		handleChampSelectEvent(data)
	default:
		fmt.Println("[EVENT DISPATCHER] unhandled event, event:" + url)
	}
}

func handleChampSelectEvent(data json.RawMessage) {
	var champSelect types.ChampSelectData
	if err := json.Unmarshal(data, &champSelect); err != nil {
		fmt.Printf("Failed to decode champion select payload: %v\n", err)
		return
	}

	var (
		wg        sync.WaitGroup
		summaries = make([]types.TeamMemberSummary, 0, len(champSelect.MyTeam))
		mu        sync.Mutex
	)

	for _, player := range champSelect.MyTeam {
		playerCopy := player

		wg.Add(1)
		go func(p types.Player) {
			defer wg.Done()

			summary := buildTeamMemberSummary(p)

			mu.Lock()
			summaries = append(summaries, summary)
			mu.Unlock()
		}(playerCopy)
	}

	wg.Wait()

	sort.SliceStable(summaries, func(i, j int) bool {
		return summaries[i].CellID < summaries[j].CellID
	})

	snapshot := types.ChampSelectSnapshot{
		QueueID:   champSelect.QueueID,
		GameID:    champSelect.GameID,
		UpdatedAt: time.Now(),
		Team:      summaries,
	}

	StoreChampSelectSnapshot(snapshot)
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
		fmt.Printf("[dispatch.buildTeamMemberSummary] failed to fetch match data: %v\n", err)
		return summary
	}

	summary.RecentMatches = buildRecentMatches(matchData.History, matchData.Heroes)

	return summary
}

func buildRecentMatches(history types.MatchHistory, heroMap map[int]controller.HeroInfo) []types.RecentMatchSummary {
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
