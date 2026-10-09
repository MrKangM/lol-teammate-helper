package dispatch

import (
	"fmt"
	"log/slog"
	"sync"

	"lol-teammate-helper/internal/analysis"
	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
	"lol-teammate-helper/internal/utils"
)

// memberInput is the lobby-side information known about one player.
type memberInput struct {
	Puuid      string
	GameName   string
	TagLine    string
	Position   string
	ChampionID int
	CellID     int
	Spells     [2]int
}

func inputFromPlayer(p types.Player) memberInput {
	return memberInput{
		Puuid: p.Puuid, GameName: p.GameName, TagLine: p.TagLine,
		Position: p.AssignedPosition, ChampionID: p.ChampionID, CellID: p.CellID,
		Spells: [2]int{p.Spell1ID, p.Spell2ID},
	}
}

func inputFromGameflow(p types.GameflowPlayer, cell int) memberInput {
	name := p.GameName
	if name == "" {
		name = p.SummonerName
	}
	return memberInput{
		Puuid: p.Puuid, GameName: name, TagLine: p.TagLine,
		Position: p.SelectedPosition, ChampionID: p.ChampionID, CellID: cell,
		Spells: [2]int{p.Spell1ID, p.Spell2ID},
	}
}

// buildMembers resolves every input concurrently and keeps the input order.
func buildMembers(inputs []memberInput) []types.TeamMemberSummary {
	out := make([]types.TeamMemberSummary, len(inputs))
	var wg sync.WaitGroup
	for i, in := range inputs {
		wg.Add(1)
		go func(i int, in memberInput) {
			defer wg.Done()
			out[i] = buildMember(in)
		}(i, in)
	}
	wg.Wait()
	return out
}

// buildMember gathers rank, recent ranked games and derived analysis for one player.
// Every remote lookup is best effort: a failure only leaves that part empty.
func buildMember(in memberInput) types.TeamMemberSummary {
	svc := matchHistorySvc

	m := types.TeamMemberSummary{
		Puuid: in.Puuid, GameName: in.GameName, TagLine: in.TagLine,
		AssignedPosition: in.Position, ChampionID: in.ChampionID, CellID: in.CellID,
		Spells: []string{service.SpellName(in.Spells[0]), service.SpellName(in.Spells[1])},
		Tags:   []string{},
	}
	if in.Puuid == "" {
		return m
	}

	var (
		wg    sync.WaitGroup
		games []types.Game
		rank  types.RankedStats
	)
	wg.Add(4)
	go func() {
		defer wg.Done()
		g, err := svc.GetRankedGames(in.Puuid, service.RankedGamesWanted)
		if err != nil {
			slog.Warn("fetch ranked games failed", "puuid", in.Puuid, "err", err)
			return
		}
		games = g
	}()
	go func() {
		defer wg.Done()
		r, err := svc.GetRankedStats(in.Puuid)
		if err != nil {
			slog.Warn("fetch rank failed", "puuid", in.Puuid, "err", err)
			return
		}
		rank = r
	}()
	go func() {
		defer wg.Done()
		m.SummonerLevel = svc.GetSummonerLevel(in.Puuid)
	}()
	go func() {
		defer wg.Done()
		m.MasteryLevel, m.MasteryPoints = svc.GetMastery(in.Puuid, in.ChampionID)
	}()
	wg.Wait()

	m.Solo = rankSummary(rank, "RANKED_SOLO_5x5")
	m.Flex = rankSummary(rank, "RANKED_FLEX_SR")

	heroIDs := []int{in.ChampionID}
	for _, g := range games {
		if len(g.Participants) > 0 {
			heroIDs = append(heroIDs, g.Participants[0].ChampionID)
		}
	}
	heroes, err := svc.GetMatchHistoryHeroesByIds(heroIDs)
	if err != nil {
		slog.Warn("fetch heroes failed", "puuid", in.Puuid, "err", err)
		heroes = map[int]types.HeroInfo{}
	}

	if hero, ok := heroes[in.ChampionID]; ok {
		m.ChampionName = hero.Name
		m.ChampionIcon = iconOf(hero)
	} else if in.ChampionID > 0 {
		m.ChampionName = fmt.Sprintf("Champion %d", in.ChampionID)
	}

	m.RecentMatches = buildRecentMatches(games, heroes)
	m.Stats = analysis.ComputeStats(games, in.ChampionID, in.Position)
	m.Rating = analysis.Rate(m)
	m.Tags = analysis.Tags(m)
	return m
}

func iconOf(h types.HeroInfo) string {
	if h.IconDataURI != "" {
		return h.IconDataURI
	}
	return h.SquarePortraitPath
}

// rankSummary extracts one queue from the raw ranked stats. Unranked queues
// come back with an empty TierKey.
func rankSummary(stats types.RankedStats, queueType string) types.RankSummary {
	rs := types.RankSummary{QueueType: queueType, QueueName: utils.QueueName(queueType)}

	entry, ok := stats.QueueMap[queueType]
	if !ok {
		for _, q := range stats.Queues {
			if q.QueueType == queueType {
				entry, ok = q, true
				break
			}
		}
	}
	if !ok || entry.Tier == "" || entry.Tier == "NA" {
		return rs
	}

	rs.TierKey = entry.Tier
	rs.Tier = utils.TierName(entry.Tier)
	rs.Division = entry.Division
	if rs.Division == "NA" {
		rs.Division = ""
	}
	rs.LeaguePoints = int(entry.LeaguePoints)
	rs.Wins = int(entry.Wins)
	rs.Losses = int(entry.Losses)
	if total := rs.Wins + rs.Losses; total > 0 {
		rs.WinRate = float64(rs.Wins) / float64(total) * 100
	}
	return rs
}

// buildRecentMatches turns the ranked games into list rows. The LCU history
// endpoint returns the queried player as the first (only) participant.
func buildRecentMatches(games []types.Game, heroes map[int]types.HeroInfo) []types.RecentMatchSummary {
	result := make([]types.RecentMatchSummary, 0, len(games))
	for _, g := range games {
		if len(g.Participants) == 0 {
			continue
		}
		p := g.Participants[0]
		st := p.Stats
		hero := heroes[p.ChampionID]

		result = append(result, types.RecentMatchSummary{
			GameID:       g.GameID,
			GameCreation: g.GameCreation,
			ChampionID:   p.ChampionID,
			ChampionName: hero.Name,
			ChampionIcon: iconOf(hero),
			Position:     analysis.Position(p.Timeline),
			Win:          st.Win,
			Kills:        st.Kills,
			Deaths:       st.Deaths,
			Assists:      st.Assists,
			CS:           st.CS(),
			Damage:       st.TotalDamageDealtToChampions,
			Gold:         st.GoldEarned,
			VisionScore:  st.VisionScore,
			QueueID:      g.QueueID,
			GameDuration: g.GameDuration,
		})
	}
	return result
}
