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
	for i := 0; i < 3; i++ {
		games = append(games, types.Game{
			GameID: int64(i), QueueID: 420,
			Participants: []types.Participant{{
				ChampionID: 1,
				Stats:      types.Stats{Win: i%2 == 0, Kills: i, Deaths: 1, Assists: 2, TotalMinionsKilled: 100, NeutralMinionsKilled: 5},
				Timeline:   types.Timeline{Lane: "MIDDLE"},
			}},
		})
	}
	games = append(games, types.Game{}) // no participants: skipped

	got := buildRecentMatches(games, map[int]types.HeroInfo{1: {Name: "Annie", SquarePortraitPath: "p"}})
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0].ChampionName != "Annie" || got[0].ChampionIcon != "p" || !got[0].Win || got[0].CS != 105 || got[0].Position != "middle" {
		t.Errorf("first = %+v", got[0])
	}
	if len(buildRecentMatches(nil, nil)) != 0 {
		t.Error("empty input should yield no rows")
	}
}

func TestRankSummary(t *testing.T) {
	stats := types.RankedStats{QueueMap: map[string]types.RankedEntry{
		"RANKED_SOLO_5x5": {Tier: "GOLD", Division: "II", LeaguePoints: 45, Wins: 30, Losses: 20},
		"RANKED_FLEX_SR":  {Tier: "NA", Division: "NA"},
	}}
	solo := rankSummary(stats, "RANKED_SOLO_5x5")
	if solo.TierKey != "GOLD" || solo.Tier != "黄金" || solo.Division != "II" || solo.WinRate != 60 || solo.QueueName != "单双排" {
		t.Errorf("solo = %+v", solo)
	}
	if flex := rankSummary(stats, "RANKED_FLEX_SR"); flex.TierKey != "" {
		t.Errorf("unranked flex should be empty, got %+v", flex)
	}
	if missing := rankSummary(types.RankedStats{}, "RANKED_SOLO_5x5"); missing.TierKey != "" {
		t.Errorf("missing queue should be empty, got %+v", missing)
	}
}

func TestContainsPuuid(t *testing.T) {
	team := []types.GameflowPlayer{{Puuid: "a"}, {Puuid: "b"}}
	if !containsPuuid(team, "b") || containsPuuid(team, "c") || containsPuuid(team, "") {
		t.Error("containsPuuid misbehaves")
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
