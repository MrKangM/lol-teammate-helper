package analysis

import (
	"reflect"
	"testing"

	"lol-teammate-helper/internal/types"
)

func game(champ int, win bool, k, d, a int, lane, role string) types.Game {
	return types.Game{Participants: []types.Participant{{
		ChampionID: champ,
		Stats:      types.Stats{Win: win, Kills: k, Deaths: d, Assists: a, TotalMinionsKilled: 100, NeutralMinionsKilled: 20, TotalDamageDealtToChampions: 20000, VisionScore: 30},
		Timeline:   types.Timeline{Lane: lane, Role: role},
	}}}
}

func TestPosition(t *testing.T) {
	cases := []struct {
		lane, role, want string
	}{
		{"TOP", "SOLO", "top"},
		{"JUNGLE", "NONE", "jungle"},
		{"MIDDLE", "SOLO", "middle"},
		{"BOTTOM", "DUO_CARRY", "bottom"},
		{"BOTTOM", "DUO_SUPPORT", "utility"},
		{"NONE", "NONE", ""},
	}
	for _, c := range cases {
		if got := Position(types.Timeline{Lane: c.lane, Role: c.role}); got != c.want {
			t.Errorf("%s/%s = %q, want %q", c.lane, c.role, got, c.want)
		}
	}
}

func TestComputeStats(t *testing.T) {
	games := []types.Game{ // newest first
		game(1, true, 10, 2, 8, "MIDDLE", "SOLO"),
		game(1, true, 5, 4, 5, "MIDDLE", "SOLO"),
		game(2, false, 2, 6, 2, "TOP", "SOLO"),
		game(1, true, 3, 3, 3, "MIDDLE", "SOLO"),
		{}, // no participants: ignored
	}
	s := ComputeStats(games, 1, "middle")

	if s.Games != 4 || s.Wins != 3 || s.WinRate != 75 {
		t.Errorf("games/wins/rate = %d/%d/%v", s.Games, s.Wins, s.WinRate)
	}
	if s.Streak != 2 {
		t.Errorf("streak = %d, want 2", s.Streak)
	}
	if s.ChampGames != 3 || s.ChampWins != 3 || s.ChampWinRate != 100 {
		t.Errorf("champ = %d/%d/%v", s.ChampGames, s.ChampWins, s.ChampWinRate)
	}
	if s.PosGames != 3 {
		t.Errorf("posGames = %d, want 3", s.PosGames)
	}
	if want := float64(20+18) / 15; s.KDA != want {
		t.Errorf("kda = %v, want %v", s.KDA, want)
	}
	if s.AvgCS != 120 {
		t.Errorf("avgCs = %v", s.AvgCS)
	}
}

func TestComputeStatsLossStreakAndEmpty(t *testing.T) {
	s := ComputeStats([]types.Game{
		game(1, false, 0, 5, 0, "", ""),
		game(1, false, 0, 5, 0, "", ""),
		game(1, true, 0, 5, 0, "", ""),
		game(1, false, 0, 5, 0, "", ""),
	}, 0, "")
	if s.Streak != -2 {
		t.Errorf("streak = %d, want -2", s.Streak)
	}
	if got := ComputeStats(nil, 1, "top"); !reflect.DeepEqual(got, types.PlayerStats{}) {
		t.Errorf("empty stats = %+v", got)
	}
}

func member(tier, div string, st types.PlayerStats, champ int) types.TeamMemberSummary {
	m := types.TeamMemberSummary{ChampionID: champ, Stats: st}
	m.Solo = types.RankSummary{TierKey: tier, Division: div}
	return m
}

func TestTierValueOrdering(t *testing.T) {
	order := [][2]string{{"IRON", "IV"}, {"IRON", "I"}, {"BRONZE", "IV"}, {"GOLD", "II"}, {"DIAMOND", "I"}, {"MASTER", ""}, {"CHALLENGER", ""}}
	prev := -1.0
	for _, o := range order {
		v, ok := TierValue(o[0], o[1])
		if !ok || v <= prev {
			t.Fatalf("%v: value %v not increasing (prev %v)", o, v, prev)
		}
		prev = v
	}
	if _, ok := TierValue("", ""); ok {
		t.Error("unranked must not have a value")
	}
}

func TestRatingLabels(t *testing.T) {
	strong := member("DIAMOND", "II", types.PlayerStats{Games: 20, Wins: 12, KDA: 4, ChampGames: 8, ChampWins: 6}, 1)
	mid := member("GOLD", "II", types.PlayerStats{Games: 20, Wins: 11, KDA: 3}, 1)
	weak := member("GOLD", "IV", types.PlayerStats{Games: 20, Wins: 8, KDA: 2}, 1)

	if r := Rate(strong); !r.Valid || r.Label != LabelTop {
		t.Errorf("strong = %+v", r)
	}
	if r := Rate(mid); !r.Valid || r.Label != LabelAvg {
		t.Errorf("mid = %+v", r)
	}
	if r := Rate(weak); !r.Valid || r.Label != LabelLow {
		t.Errorf("weak = %+v", r)
	}
}

func TestRatingNeedsData(t *testing.T) {
	if r := Rate(types.TeamMemberSummary{}); r.Valid || r.Label != "" {
		t.Errorf("empty member rated: %+v", r)
	}
	// rank only is enough for a provisional verdict
	if r := Rate(member("PLATINUM", "I", types.PlayerStats{}, 0)); !r.Valid {
		t.Error("rank-only member should be rated")
	}
}

func TestTags(t *testing.T) {
	m := types.TeamMemberSummary{
		ChampionID: 1, AssignedPosition: "middle", SummonerLevel: 30,
		Stats: types.PlayerStats{Games: 20, Streak: -4, ChampGames: 0, PosGames: 0, AvgDeaths: 9},
	}
	got := Tags(m)
	want := []string{"4连败", "新英雄", "非常用位置", "高死亡", "低等级"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tags = %v, want %v", got, want)
	}

	m = types.TeamMemberSummary{ChampionID: 1, Stats: types.PlayerStats{Games: 20, Streak: 5, ChampGames: 6, ChampWinRate: 70, PosGames: 10}}
	if got := Tags(m); !reflect.DeepEqual(got, []string{"5连胜", "本命英雄"}) {
		t.Errorf("tags = %v", got)
	}
}
