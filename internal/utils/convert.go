package utils

import (
	"fmt"
	"lol-teammate-helper/internal/types"
	"strings"
)

func ConvertRankDataToChinese(rankInfo *types.RankedStats) {
	tierMap := map[string]string{
		"IRON":        "\u9ed1\u94c1",
		"BRONZE":      "\u9752\u94dc",
		"SILVER":      "\u767d\u94f6",
		"GOLD":        "\u9ec4\u91d1",
		"PLATINUM":    "\u767d\u91d1",
		"EMERALD":     "\u7fe1\u7fe0",
		"DIAMOND":     "\u94bb\u77f3",
		"MASTER":      "\u5927\u5e08",
		"GRANDMASTER": "\u5b97\u5e08",
		"CHALLENGER":  "\u738b\u8005",
		"":            "\u672a\u5b9a\u7ea7",
		"NA":          "\u672a\u5b9a\u7ea7",
	}

	queueMap := map[string]string{
		"RANKED_SOLO_5x5":      "\u5355\u53cc\u6392",
		"RANKED_FLEX_SR":       "\u7075\u6d3b\u7ec4\u6392",
		"RANKED_TFT":           "\u4e91\u9876\u4e4b\u5f08",
		"RANKED_TFT_DOUBLE_UP": "\u4e91\u9876\u53cc\u4eba",
		"RANKED_TFT_TURBO":     "\u4e91\u9876\u6781\u901f",
	}

	if chineseTier, ok := tierMap[rankInfo.HighestCurrentSeasonReachedTierSR]; ok {
		rankInfo.HighestCurrentSeasonReachedTierSR = chineseTier
	}
	if chineseTier, ok := tierMap[rankInfo.HighestPreviousSeasonEndTier]; ok {
		rankInfo.HighestPreviousSeasonEndTier = chineseTier
	}

	for queueType, entry := range rankInfo.QueueMap {
		convertRankedEntry(&entry, tierMap, queueMap, queueType)
		rankInfo.QueueMap[queueType] = entry
	}

	for i := range rankInfo.Queues {
		convertRankedEntry(&rankInfo.Queues[i], tierMap, queueMap, rankInfo.Queues[i].QueueType)
	}

	convertRankedEntry(&rankInfo.HighestRankedEntry, tierMap, queueMap, rankInfo.HighestRankedEntry.QueueType)
	convertRankedEntry(&rankInfo.HighestRankedEntrySR, tierMap, queueMap, rankInfo.HighestRankedEntrySR.QueueType)
}

func convertRankedEntry(entry *types.RankedEntry, tierMap map[string]string, queueMap map[string]string, queueType string) {
	if chineseTier, ok := tierMap[entry.Tier]; ok {
		entry.Tier = chineseTier
	}
	if chineseTier, ok := tierMap[entry.HighestTier]; ok {
		entry.HighestTier = chineseTier
	}
	if chineseTier, ok := tierMap[entry.PreviousSeasonEndTier]; ok {
		entry.PreviousSeasonEndTier = chineseTier
	}
	if chineseTier, ok := tierMap[entry.PreviousSeasonHighestTier]; ok {
		entry.PreviousSeasonHighestTier = chineseTier
	}
	if chineseTier, ok := tierMap[entry.RatedTier]; ok {
		entry.RatedTier = chineseTier
	}

	if chineseName, ok := queueMap[queueType]; ok {
		entry.QueueDisplayName = chineseName
	}
}

func FormatRankInfo(rankInfo types.RankedStats) string {
	var result strings.Builder

	result.WriteString("=== \u6392\u4f4d\u4fe1\u606f\u6c47\u603b ===\n")
	result.WriteString(fmt.Sprintf("\u5f53\u524d\u8d5b\u5b63\u6700\u9ad8\u6bb5\u4f4d: %s\n", rankInfo.HighestCurrentSeasonReachedTierSR))
	result.WriteString(fmt.Sprintf("\u4e0a\u8d5b\u5b63\u6700\u9ad8\u6bb5\u4f4d: %s %s\n",
		rankInfo.HighestPreviousSeasonEndTier,
		rankInfo.HighestPreviousSeasonEndDivision))

	for queueType, entry := range rankInfo.QueueMap {
		if entry.Tier == "" || entry.Tier == "\u672a\u5b9a\u7ea7" {
			continue
		}

		queueName := entry.QueueDisplayName
		if queueName == "" {
			queueName = getQueueChineseName(queueType)
		}

		result.WriteString(fmt.Sprintf("\n[%s]\n", queueName))
		result.WriteString(fmt.Sprintf("\u5f53\u524d\u6bb5\u4f4d: %s %s", entry.Tier, entry.Division))
		if entry.LeaguePoints > 0 {
			result.WriteString(fmt.Sprintf(" (%d\u80dc\u70b9)", entry.LeaguePoints))
		}
		result.WriteString("\n")

		if entry.HighestTier != "" && entry.HighestTier != "\u672a\u5b9a\u7ea7" {
			result.WriteString(fmt.Sprintf("\u5386\u53f2\u6700\u9ad8: %s %s\n", entry.HighestTier, entry.HighestDivision))
		}

		if total := entry.Wins + entry.Losses; total > 0 {
			winRate := float64(entry.Wins) / float64(total) * 100
			result.WriteString(fmt.Sprintf("\u6218\u7ee9: %d\u80dc %d\u8d1f (\u80dc\u7387%.1f%%)\n", entry.Wins, entry.Losses, winRate))
		}
	}

	return result.String()
}

func getQueueChineseName(queueType string) string {
	queueNames := map[string]string{
		"RANKED_SOLO_5x5":      "\u5355\u53cc\u6392",
		"RANKED_FLEX_SR":       "\u7075\u6d3b\u7ec4\u6392",
		"RANKED_TFT":           "\u4e91\u9876\u4e4b\u5f08",
		"RANKED_TFT_DOUBLE_UP": "\u4e91\u9876\u53cc\u4eba",
		"RANKED_TFT_TURBO":     "\u4e91\u9876\u6781\u901f",
	}
	if name, ok := queueNames[queueType]; ok {
		return name
	}
	return queueType
}
