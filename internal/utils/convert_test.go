package utils

import (
	"testing"

	"lol-teammate-helper/internal/types"
)

func TestConvertRankDataToChinese(t *testing.T) {
	stats := types.RankedStats{
		HighestCurrentSeasonReachedTierSR: "GOLD",
		QueueMap: map[string]types.RankedEntry{
			"RANKED_SOLO_5x5": {QueueType: "RANKED_SOLO_5x5", Tier: "DIAMOND", HighestTier: "MASTER"},
		},
		Queues: []types.RankedEntry{{QueueType: "RANKED_FLEX_SR", Tier: ""}},
	}
	ConvertRankDataToChinese(&stats)

	if stats.HighestCurrentSeasonReachedTierSR != "黄金" {
		t.Errorf("tier = %q", stats.HighestCurrentSeasonReachedTierSR)
	}
	solo := stats.QueueMap["RANKED_SOLO_5x5"]
	if solo.Tier != "钻石" || solo.HighestTier != "大师" || solo.QueueDisplayName != "单双排" {
		t.Errorf("solo = %+v", solo)
	}
	if stats.Queues[0].Tier != "未定级" || stats.Queues[0].QueueDisplayName != "灵活组排" {
		t.Errorf("flex = %+v", stats.Queues[0])
	}
}

func TestGetServerChineseName(t *testing.T) {
	if got := GetServerChineseName("HN1"); got != "艾欧尼亚" {
		t.Errorf("got %q", got)
	}
	if got := GetServerChineseName("UNKNOWN"); got != "UNKNOWN" {
		t.Errorf("got %q", got)
	}
}
