package dispatch

import (
	"encoding/json"
	"os"
	"testing"

	"lol-teammate-helper/internal/types"
)

func loadFixture(t *testing.T) types.WSMessageType {
	t.Helper()
	raw, err := os.ReadFile("testdata/champ_select_session.json")
	if err != nil {
		t.Fatal(err)
	}
	var msg types.WSMessageType
	if err := json.Unmarshal(raw, &msg); err != nil {
		t.Fatal(err)
	}
	return msg
}

func TestFixtureDecodes(t *testing.T) {
	msg := loadFixture(t)
	var cs types.ChampSelectData
	if err := json.Unmarshal(msg.Data, &cs); err != nil {
		t.Fatal(err)
	}
	if len(cs.MyTeam) != 5 || cs.QueueID != 420 {
		t.Fatalf("unexpected fixture content: %d players, queue %d", len(cs.MyTeam), cs.QueueID)
	}
}

func TestBuildRecentMatches(t *testing.T) {
	var games []types.Game
	for i := 0; i < 8; i++ {
		games = append(games, types.Game{
			QueueID: 420,
			Participants: []types.Participant{{
				ChampionID: 1,
				Stats:      types.Stats{Win: i%2 == 0, Kills: i, Deaths: 1, Assists: 2},
			}},
		})
	}
	games = append(games[:2], append([]types.Game{{}}, games[2:]...)...) // a game with no participants is skipped

	history := types.MatchHistory{Games: types.MatchHistoryGames{Games: games}}
	got := buildRecentMatches(history, map[int]types.HeroInfo{1: {Name: "Annie", SquarePortraitPath: "p"}})

	if len(got) != maxRecentMatches {
		t.Fatalf("len = %d, want %d", len(got), maxRecentMatches)
	}
	if got[0].ChampionName != "Annie" || got[0].ChampionIcon != "p" || !got[0].Win {
		t.Errorf("first = %+v", got[0])
	}
	if buildRecentMatches(types.MatchHistory{}, nil) != nil {
		t.Error("empty history should yield nil")
	}
}

func TestStoreSnapshotDeduplicates(t *testing.T) {
	ClearChampSelectSnapshot()
	snap := types.ChampSelectSnapshot{QueueID: 420, Team: []types.TeamMemberSummary{{Puuid: "a"}}}

	if !StoreChampSelectSnapshot(snap) {
		t.Fatal("first store should report a change")
	}
	same := snap
	same.UpdatedAt = same.UpdatedAt.Add(1)
	if StoreChampSelectSnapshot(same) {
		t.Fatal("identical content should not report a change")
	}
	snap.Team = []types.TeamMemberSummary{{Puuid: "a", ChampionID: 3}}
	if !StoreChampSelectSnapshot(snap) {
		t.Fatal("changed content should report a change")
	}

	ClearChampSelectSnapshot()
	if _, ok := GetChampSelectSnapshot(); ok {
		t.Fatal("snapshot should be cleared")
	}
}
