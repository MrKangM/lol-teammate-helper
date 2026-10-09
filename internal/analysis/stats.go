// Package analysis turns raw match history into the statistics, rating and
// tags shown on the teammate cards. Everything here is pure and unit tested.
package analysis

import (
	"strings"

	"lol-teammate-helper/internal/types"
)

// Ranked queue IDs the app cares about.
const (
	QueueSolo = 420
	QueueFlex = 440
)

// IsRankedQueue reports whether the queue is ranked solo/duo or flex.
func IsRankedQueue(queueID int) bool {
	return queueID == QueueSolo || queueID == QueueFlex
}

// Position maps the LCU lane/role pair to a champ-select style position name
// (top / jungle / middle / bottom / utility). Unknown combinations yield "".
func Position(tl types.Timeline) string {
	lane := strings.ToUpper(tl.Lane)
	role := strings.ToUpper(tl.Role)
	switch lane {
	case "TOP":
		return "top"
	case "JUNGLE":
		return "jungle"
	case "MID", "MIDDLE":
		return "middle"
	case "BOTTOM", "BOT":
		if role == "DUO_SUPPORT" || role == "SUPPORT" {
			return "utility"
		}
		return "bottom"
	}
	return ""
}

// ComputeStats aggregates games (newest first, each holding the player as
// Participants[0]). championID and position (both optional) enable the
// current-champion and current-position breakdown.
func ComputeStats(games []types.Game, championID int, position string) types.PlayerStats {
	var s types.PlayerStats
	var kills, deaths, assists, cs, damage, vision int

	streakOpen := true
	for _, g := range games {
		if len(g.Participants) == 0 {
			continue
		}
		p := g.Participants[0]
		st := p.Stats

		s.Games++
		kills += st.Kills
		deaths += st.Deaths
		assists += st.Assists
		cs += st.CS()
		damage += st.TotalDamageDealtToChampions
		vision += st.VisionScore
		if st.Win {
			s.Wins++
		}

		if streakOpen {
			switch {
			case s.Streak == 0 || (s.Streak > 0) == st.Win:
				if st.Win {
					s.Streak++
				} else {
					s.Streak--
				}
			default:
				streakOpen = false
			}
		}

		if championID > 0 && p.ChampionID == championID {
			s.ChampGames++
			if st.Win {
				s.ChampWins++
			}
		}
		if position != "" && Position(p.Timeline) == position {
			s.PosGames++
		}
	}

	if s.Games == 0 {
		return s
	}
	n := float64(s.Games)
	s.WinRate = pct(s.Wins, s.Games)
	s.AvgKills = float64(kills) / n
	s.AvgDeaths = float64(deaths) / n
	s.AvgAssists = float64(assists) / n
	s.KDA = float64(kills+assists) / maxf(float64(deaths), 1)
	s.AvgCS = float64(cs) / n
	s.AvgDamage = float64(damage) / n
	s.AvgVision = float64(vision) / n
	if s.ChampGames > 0 {
		s.ChampWinRate = pct(s.ChampWins, s.ChampGames)
	}
	return s
}

func pct(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b) * 100
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
