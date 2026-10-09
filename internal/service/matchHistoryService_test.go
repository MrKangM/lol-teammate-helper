package service

import (
	"testing"
	"time"

	"lol-teammate-helper/internal/types"
)

func TestCollectChampionIDs(t *testing.T) {
	h := types.MatchHistory{Games: types.MatchHistoryGames{Games: []types.Game{
		{Participants: []types.Participant{{ChampionID: 1}, {ChampionID: 0}}},
		{Participants: []types.Participant{{ChampionID: 1}, {ChampionID: 2}}},
	}}}
	if got := collectChampionIDs(h); len(got) != 2 {
		t.Fatalf("got %v, want 2 unique positive ids", got)
	}
	if collectChampionIDs(types.MatchHistory{}) != nil {
		t.Fatal("empty history should yield nil")
	}
}

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
