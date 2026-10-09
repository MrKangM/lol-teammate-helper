package service

import "testing"

func TestAllowedAssetPath(t *testing.T) {
	ok := []string{
		"/lol-game-data/assets/v1/profile-icons/1.jpg",
		"/fe/lol-static-assets/images/ranked-emblem/emblem-gold.png",
		"/lol-static-assets/v1/x.png",
	}
	bad := []string{
		"", "https://evil.example/x.png", "/lol-summoner/v1/current-summoner",
		"/lol-game-data/../lol-chat/v1/me", "/fe/..\\x", "//evil.example/x", "/fe/x\r\nHost: y",
	}
	for _, p := range ok {
		if !AllowedAssetPath(p) {
			t.Errorf("%q should be allowed", p)
		}
	}
	for _, p := range bad {
		if AllowedAssetPath(p) {
			t.Errorf("%q should be rejected", p)
		}
	}
}

func TestRankEmblemCandidates(t *testing.T) {
	c := RankEmblemCandidates("GOLD")
	if len(c) == 0 || c[0] != "/fe/lol-static-assets/images/ranked-emblem/emblem-gold.png" {
		t.Fatalf("got %v", c)
	}
	if got := RankEmblemCandidates("")[0]; got != "/fe/lol-static-assets/images/ranked-emblem/emblem-unranked.png" {
		t.Errorf("unranked = %q", got)
	}
	for _, p := range append(c, PositionIconCandidates("utility")...) {
		if !AllowedAssetPath(p) {
			t.Errorf("built-in candidate %q is not allowed", p)
		}
	}
}

func TestGetAssetWithoutClientIsEmptyAndRejectsBadPaths(t *testing.T) {
	svc := NewMatchHistoryService()
	if got := svc.GetAsset([]string{"/lol-game-data/x.png"}); got != "" {
		t.Errorf("expected empty without credentials, got %q", got)
	}
}
