package service

import (
	"testing"
	"time"

	"lol-teammate-helper/internal/types"
)

func TestMatchCacheExpires(t *testing.T) {
	svc := NewMatchHistoryService()
	clock := time.Unix(1000, 0)
	svc.now = func() time.Time { return clock }

	svc.matchCache["p"] = matchCacheEntry{history: types.MatchHistory{AccountID: 7}, expiresAt: clock.Add(matchCacheTTL)}
	if h, ok := svc.getMatchesFromCache("p"); !ok || h.AccountID != 7 {
		t.Fatal("expected cache hit")
	}
	clock = clock.Add(matchCacheTTL + time.Second)
	if _, ok := svc.getMatchesFromCache("p"); ok {
		t.Fatal("expected cache miss after TTL")
	}
}

func TestHeroesServedFromCache(t *testing.T) {
	svc := NewMatchHistoryService()
	svc.storeHeroInCache(5, types.HeroInfo{Name: "X"})
	got, err := svc.GetMatchHistoryHeroesByIds([]int{5, 5, -1})
	if err != nil || len(got) != 1 || got[5].Name != "X" {
		t.Fatalf("got %v, err %v", got, err)
	}
}

func TestBuildGameDetail(t *testing.T) {
	p := func(id, team, champ int, win bool, kills int, item0 int) types.Participant {
		return types.Participant{ParticipantID: id, TeamID: team, ChampionID: champ, Spell1ID: 4, Spell2ID: 14,
			Stats: types.Stats{Win: win, Kills: kills, GoldEarned: 1000, Item0: item0, TotalMinionsKilled: 50, NeutralMinionsKilled: 10}}
	}
	game := types.Game{
		GameID: 9, QueueID: 420, GameDuration: 1800,
		Participants: []types.Participant{p(1, 100, 1, true, 5, 3000), p(2, 100, 2, true, 3, 0), p(3, 200, 3, false, 4, 0)},
		ParticipantIdentities: []types.ParticipantIdentity{
			{ParticipantID: 1, Player: types.IdentityPlayer{Puuid: "me", GameName: "Me", TagLine: "123"}},
			{ParticipantID: 2, Player: types.IdentityPlayer{Puuid: "ally", SummonerName: "Ally"}},
			{ParticipantID: 3, Player: types.IdentityPlayer{Puuid: "foe", GameName: "Foe"}},
		},
		Teams: []types.GameTeam{{TeamID: 100, Win: "Win"}, {TeamID: 200, Win: "Fail"}},
	}
	heroes := map[int]types.HeroInfo{1: {Name: "Annie", IconDataURI: "data:a"}, 2: {Name: "Olaf", SquarePortraitPath: "p"}}
	d := BuildGameDetail(game, "me", heroes, map[int]string{3000: "data:item"}, SpellName)

	if len(d.Teams) != 2 || !d.Teams[0].Win || d.Teams[1].Win {
		t.Fatalf("teams = %+v", d.Teams)
	}
	blue := d.Teams[0]
	if blue.Kills != 8 || blue.Gold != 2000 || len(blue.Players) != 2 {
		t.Errorf("blue = %+v", blue)
	}
	me := blue.Players[0]
	if me.Name != "Me #123" || !me.IsTarget || me.ChampionIcon != "data:a" || me.CS != 60 || me.Items[0] != "data:item" || len(me.Items) != 7 {
		t.Errorf("me = %+v", me)
	}
	if me.Spells[0] != "闪现" || me.Spells[1] != "引燃" {
		t.Errorf("spells = %v", me.Spells)
	}
	ally := blue.Players[1]
	if ally.Name != "Ally" || ally.IsTarget || ally.ChampionIcon != "p" {
		t.Errorf("ally = %+v", ally)
	}
}
