package dispatch

import (
	"encoding/json"
	"log/slog"

	"lol-teammate-helper/internal/diag"
	"lol-teammate-helper/internal/types"
)

// handleGameflowEvent reacts to the game phase. Once the game is loading or
// running the LCU lists all ten players, which is the first moment the
// opposing team's identities are available.
func handleGameflowEvent(data json.RawMessage) {
	var session types.GameflowSession
	if err := json.Unmarshal(data, &session); err != nil {
		slog.Warn("decode gameflow payload failed", "err", err)
		return
	}

	diag.SetPhase(session.Phase)
	slog.Info("gameflow phase", "phase", session.Phase,
		"teamOne", len(session.GameData.TeamOne), "teamTwo", len(session.GameData.TeamTwo))

	switch session.Phase {
	case "GameStart", "InProgress", "Reconnect":
		publishInGame(session)
	case "None", "Lobby", "Matchmaking", "ReadyCheck":
		handleSessionEnded()
	}
	// ChampSelect is handled by the champ-select topic; end-of-game phases keep the last view.
}

func publishInGame(session types.GameflowSession) {
	me := matchHistorySvc.GetCurrentPuuid()
	one, two := session.GameData.TeamOne, session.GameData.TeamTwo

	mine, theirs := one, two
	if containsPuuid(two, me) {
		mine, theirs = two, one
	} else if !containsPuuid(one, me) {
		slog.Debug("current player not found in gameflow teams; assuming team one")
	}

	toInputs := func(players []types.GameflowPlayer) []memberInput {
		var in []memberInput
		for i, p := range players {
			if p.Puuid != "" {
				in = append(in, inputFromGameflow(p, i))
			}
		}
		return in
	}
	ownInputs, enemyInputs := toInputs(mine), toInputs(theirs)
	if len(ownInputs) == 0 && len(enemyInputs) == 0 {
		return
	}

	publish(types.ChampSelectSnapshot{
		QueueID: session.GameData.Queue.ID,
		GameID:  session.GameData.GameID,
		Phase:   session.Phase,
		Team:    buildMembers(ownInputs),
		Enemy:   buildMembers(enemyInputs),
	})
}

func containsPuuid(players []types.GameflowPlayer, puuid string) bool {
	if puuid == "" {
		return false
	}
	for _, p := range players {
		if p.Puuid == puuid {
			return true
		}
	}
	return false
}
