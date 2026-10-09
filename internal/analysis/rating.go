package analysis

import (
	"strconv"

	"lol-teammate-helper/internal/types"
)

// Label thresholds on the 0-100 composite score. Tweak here to retune the tool.
const (
	thresholdTop   = 68.0 // 大腿
	thresholdGood  = 55.0 // 上等马
	thresholdAvg   = 42.0 // 中等马; below this is 下等马
	minGamesRating = 5    // below this many recent ranked games we have no verdict

	// Component weights; missing components are dropped and the rest renormalised.
	weightWinRate = 0.35
	weightKDA     = 0.25
	weightRank    = 0.25
	weightChamp   = 0.15
)

const (
	LabelTop  = "大腿"
	LabelGood = "上等马"
	LabelAvg  = "中等马"
	LabelLow  = "下等马"
)

var tierIndex = map[string]int{
	"IRON": 0, "BRONZE": 1, "SILVER": 2, "GOLD": 3, "PLATINUM": 4,
	"EMERALD": 5, "DIAMOND": 6, "MASTER": 7, "GRANDMASTER": 8, "CHALLENGER": 9,
}

var divisionIndex = map[string]int{"IV": 0, "III": 1, "II": 2, "I": 3}

// TierValue places a rank on a linear ladder: Iron IV = 0, +1 per division,
// Master+ occupy the slots after Diamond I. Unranked returns ok=false.
func TierValue(tierKey, division string) (float64, bool) {
	idx, ok := tierIndex[tierKey]
	if !ok {
		return 0, false
	}
	if idx >= tierIndex["MASTER"] {
		return float64(tierIndex["MASTER"]*4) + float64(idx-tierIndex["MASTER"]), true
	}
	return float64(idx*4 + divisionIndex[division]), true
}

// BestRank prefers the solo/duo entry and falls back to flex.
func BestRank(m types.TeamMemberSummary) (types.RankSummary, bool) {
	if m.Solo.TierKey != "" {
		return m.Solo, true
	}
	if m.Flex.TierKey != "" {
		return m.Flex, true
	}
	return types.RankSummary{}, false
}

// Score computes the composite 0-100 strength estimate for a player from
// recent ranked performance, rank and familiarity with the current champion.
// It returns ok=false when there is too little data to say anything.
func Score(m types.TeamMemberSummary) (float64, bool) {
	st := m.Stats

	type part struct{ value, weight float64 }
	var parts []part

	if st.Games >= minGamesRating {
		// Smooth with 5 phantom 50% games so a 3-0 start doesn't read as 100%.
		wr := (float64(st.Wins) + 2.5) / (float64(st.Games) + 5)
		parts = append(parts,
			part{scale(wr, 0.35, 0.65), weightWinRate},
			part{scale(st.KDA, 1.5, 4.5), weightKDA},
		)
	}

	if rk, ok := BestRank(m); ok {
		if v, ok := TierValue(rk.TierKey, rk.Division); ok {
			parts = append(parts, part{scale(v, 4, 30), weightRank})
		}
	}

	if m.ChampionID > 0 && st.Games >= minGamesRating {
		champ := 35.0 // never played it recently
		if st.ChampGames > 0 {
			wr := (float64(st.ChampWins) + 1) / (float64(st.ChampGames) + 2)
			familiarity := minf(float64(st.ChampGames)/8, 1)
			champ = 35 + 65*familiarity*scale(wr, 0.3, 0.7)/100
		}
		parts = append(parts, part{champ, weightChamp})
	}

	var total, weights float64
	for _, p := range parts {
		total += p.value * p.weight
		weights += p.weight
	}
	// Without recent games we only trust the rank, and only if it is present.
	if weights == 0 || (st.Games < minGamesRating && weights < weightRank) {
		return 0, false
	}
	return total / weights, true
}

// LabelFor maps a score to the horse label.
func LabelFor(score float64) string {
	switch {
	case score >= thresholdTop:
		return LabelTop
	case score >= thresholdGood:
		return LabelGood
	case score >= thresholdAvg:
		return LabelAvg
	default:
		return LabelLow
	}
}

// Rate fills in the Rating of a member.
func Rate(m types.TeamMemberSummary) types.Rating {
	score, ok := Score(m)
	if !ok {
		return types.Rating{}
	}
	return types.Rating{Score: score, Valid: true, Label: LabelFor(score)}
}

// Tags derives short human-readable hints from a member's stats.
func Tags(m types.TeamMemberSummary) []string {
	var tags []string
	st := m.Stats

	switch {
	case st.Streak >= 3:
		tags = append(tags, strconv.Itoa(st.Streak)+"连胜")
	case st.Streak <= -3:
		tags = append(tags, strconv.Itoa(-st.Streak)+"连败")
	}

	if m.ChampionID > 0 && st.Games >= 10 {
		switch {
		case st.ChampGames == 0:
			tags = append(tags, "新英雄")
		case st.ChampGames >= 5 && st.ChampWinRate >= 55:
			tags = append(tags, "本命英雄")
		case st.ChampGames >= 5 && st.ChampWinRate < 40:
			tags = append(tags, "英雄胜率低")
		}
	}

	if m.AssignedPosition != "" && st.Games >= 10 && st.PosGames <= 1 {
		tags = append(tags, "非常用位置")
	}
	if st.Games >= minGamesRating && st.AvgDeaths >= 8 {
		tags = append(tags, "高死亡")
	}
	if m.SummonerLevel > 0 && m.SummonerLevel < 40 {
		tags = append(tags, "低等级")
	}
	return tags
}

// scale maps v from [lo, hi] onto [0, 100], clamped.
func scale(v, lo, hi float64) float64 {
	if hi <= lo {
		return 0
	}
	r := (v - lo) / (hi - lo) * 100
	if r < 0 {
		return 0
	}
	if r > 100 {
		return 100
	}
	return r
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
