package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"lol-teammate-helper/internal/analysis"
	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/types"
)

const (
	rankedStatsPath = "/lol-ranked/v1/ranked-stats/%s"
	summonerPath    = "/lol-summoner/v2/summoners/puuid/%s"
	masteryPath     = "/lol-champion-mastery/v1/%s/champion-mastery"
	matchesPagePath = "/lol-match-history/v1/products/lol/%s/matches?begIndex=%d&endIndex=%d"
	gameDetailPath  = "/lol-match-history/v1/games/%d"
	itemsPath       = "/lol-game-data/assets/v1/items.json"

	// RankedGamesWanted is how many recent ranked games are analysed per player.
	RankedGamesWanted = 20

	rankedPageSize = 40
	rankedMaxPages = 3

	slowCacheTTL = 10 * time.Minute // rank, level and mastery change rarely
)

type bodyCacheEntry struct {
	body      []byte
	expiresAt time.Time
}

type gamesCacheEntry struct {
	games     []types.Game
	expiresAt time.Time
}

// cachedGet fetches an LCU endpoint, caching the raw body for ttl.
// Callers unmarshal the body themselves so cached state is never shared.
func (svc *MatchHistoryService) cachedGet(endpoint string, ttl time.Duration) ([]byte, error) {
	svc.cacheMu.RLock()
	entry, ok := svc.bodyCache[endpoint]
	svc.cacheMu.RUnlock()
	if ok && svc.now().Before(entry.expiresAt) {
		return entry.body, nil
	}

	cfg, ok := config.Instance()
	if !ok {
		return nil, errors.New("riot credentials are not initialised")
	}
	body, err := cfg.SendHttpRequest(endpoint, http.MethodGet)
	if err != nil {
		return nil, err
	}

	svc.cacheMu.Lock()
	svc.bodyCache[endpoint] = bodyCacheEntry{body: body, expiresAt: svc.now().Add(ttl)}
	svc.cacheMu.Unlock()
	return body, nil
}

// GetRankedStats returns the raw (untranslated) ranked stats of a player.
func (svc *MatchHistoryService) GetRankedStats(puuid string) (types.RankedStats, error) {
	var stats types.RankedStats
	body, err := svc.cachedGet(fmt.Sprintf(rankedStatsPath, puuid), slowCacheTTL)
	if err != nil {
		return stats, err
	}
	err = json.Unmarshal(body, &stats)
	return stats, err
}

// GetSummonerLevel returns the account level, or 0 when it cannot be fetched.
func (svc *MatchHistoryService) GetSummonerLevel(puuid string) int {
	body, err := svc.cachedGet(fmt.Sprintf(summonerPath, puuid), slowCacheTTL)
	if err != nil {
		return 0
	}
	var s struct {
		SummonerLevel int `json:"summonerLevel"`
	}
	if json.Unmarshal(body, &s) != nil {
		return 0
	}
	return s.SummonerLevel
}

// GetMastery returns the mastery level and points for one champion (0, 0 if unknown).
func (svc *MatchHistoryService) GetMastery(puuid string, championID int) (level, points int) {
	if championID <= 0 {
		return 0, 0
	}
	body, err := svc.cachedGet(fmt.Sprintf(masteryPath, puuid), slowCacheTTL)
	if err != nil {
		return 0, 0
	}
	var list []struct {
		ChampionID     int `json:"championId"`
		ChampionLevel  int `json:"championLevel"`
		ChampionPoints int `json:"championPoints"`
	}
	if json.Unmarshal(body, &list) != nil {
		return 0, 0
	}
	for _, m := range list {
		if m.ChampionID == championID {
			return m.ChampionLevel, m.ChampionPoints
		}
	}
	return 0, 0
}

// GetRankedGames returns up to want most recent solo/duo and flex games
// (newest first). The history endpoint mixes all queues, so it pages until
// enough ranked games are found or the page budget runs out.
func (svc *MatchHistoryService) GetRankedGames(puuid string, want int) ([]types.Game, error) {
	svc.cacheMu.RLock()
	entry, ok := svc.gamesCache[puuid]
	svc.cacheMu.RUnlock()
	if ok && svc.now().Before(entry.expiresAt) && len(entry.games) >= want {
		return entry.games[:want], nil
	}

	cfg, ok := config.Instance()
	if !ok {
		return nil, errors.New("riot credentials are not initialised")
	}

	var games []types.Game
	for page := 0; page < rankedMaxPages && len(games) < want; page++ {
		begin := page * rankedPageSize
		body, err := cfg.SendHttpRequest(fmt.Sprintf(matchesPagePath, puuid, begin, begin+rankedPageSize), http.MethodGet)
		if err != nil {
			if page == 0 {
				return nil, err
			}
			break // keep what the earlier pages gave us
		}
		var history types.MatchHistory
		if err := json.Unmarshal(body, &history); err != nil {
			return nil, err
		}
		for _, g := range history.Games.Games {
			if analysis.IsRankedQueue(g.QueueID) {
				games = append(games, g)
			}
		}
		if len(history.Games.Games) < rankedPageSize {
			break // reached the end of the history
		}
	}
	if len(games) > want {
		games = games[:want]
	}

	svc.cacheMu.Lock()
	svc.gamesCache[puuid] = gamesCacheEntry{games: games, expiresAt: svc.now().Add(matchCacheTTL)}
	svc.cacheMu.Unlock()
	return games, nil
}

// GetGameDetail builds the end-of-game scoreboard of one match.
// highlightPuuid marks the player the user clicked on (optional).
func (svc *MatchHistoryService) GetGameDetail(gameID int64, highlightPuuid string) (types.GameDetail, error) {
	body, err := svc.cachedGet(fmt.Sprintf(gameDetailPath, gameID), slowCacheTTL)
	if err != nil {
		return types.GameDetail{}, err
	}
	var game types.Game
	if err := json.Unmarshal(body, &game); err != nil {
		return types.GameDetail{}, err
	}

	var heroIDs, itemIDs []int
	for _, p := range game.Participants {
		heroIDs = append(heroIDs, p.ChampionID)
		for _, id := range p.Stats.Items() {
			if id > 0 {
				itemIDs = append(itemIDs, id)
			}
		}
	}

	heroes, err := svc.GetMatchHistoryHeroesByIds(heroIDs)
	if err != nil {
		slog.Warn("load heroes for game detail failed", "game", gameID, "err", err)
		heroes = map[int]types.HeroInfo{}
	}
	icons := svc.GetItemIcons(itemIDs)

	return BuildGameDetail(game, highlightPuuid, heroes, icons, svc.GetSpell), nil
}

// BuildGameDetail assembles the scoreboard DTO from raw game data; it has no
// side effects so it can be unit tested.
func BuildGameDetail(game types.Game, highlightPuuid string, heroes map[int]types.HeroInfo, itemIcons map[int]string, spell func(int) types.SpellInfo) types.GameDetail {
	identities := make(map[int]types.IdentityPlayer, len(game.ParticipantIdentities))
	for _, id := range game.ParticipantIdentities {
		identities[id.ParticipantID] = id.Player
	}
	teamWin := make(map[int]bool, len(game.Teams))
	for _, t := range game.Teams {
		teamWin[t.TeamID] = t.Win == "Win"
	}

	detail := types.GameDetail{
		GameID:       game.GameID,
		GameCreation: game.GameCreation,
		GameDuration: game.GameDuration,
		QueueID:      game.QueueID,
	}
	byTeam := map[int]*types.GameDetailTeam{}

	for _, p := range game.Participants {
		team, ok := byTeam[p.TeamID]
		if !ok {
			detail.Teams = append(detail.Teams, types.GameDetailTeam{TeamID: p.TeamID, Win: p.Stats.Win})
			team = &detail.Teams[len(detail.Teams)-1]
			byTeam[p.TeamID] = team
			if w, known := teamWin[p.TeamID]; known {
				team.Win = w
			}
		}

		who := identities[p.ParticipantID]
		name := who.GameName
		if name == "" {
			name = who.SummonerName
		}
		if who.TagLine != "" && who.GameName != "" {
			name += " #" + who.TagLine
		}

		hero := heroes[p.ChampionID]
		icon := hero.IconDataURI
		if icon == "" {
			icon = hero.SquarePortraitPath
		}

		var items []string
		for _, id := range p.Stats.Items() {
			items = append(items, itemIcons[id]) // "" for empty or unresolved slots
		}

		st := p.Stats
		team.Kills += st.Kills
		team.Gold += st.GoldEarned
		team.Players = append(team.Players, types.GameDetailPlayer{
			Puuid:        who.Puuid,
			Name:         name,
			ChampionID:   p.ChampionID,
			ChampionName: hero.Name,
			ChampionIcon: icon,
			Level:        st.ChampLevel,
			Kills:        st.Kills,
			Deaths:       st.Deaths,
			Assists:      st.Assists,
			CS:           st.CS(),
			Gold:         st.GoldEarned,
			Damage:       st.TotalDamageDealtToChampions,
			DamageTaken:  st.TotalDamageTaken,
			VisionScore:  st.VisionScore,
			Spells:       []types.SpellInfo{spell(p.Spell1ID), spell(p.Spell2ID)},
			Items:        items,
			IsTarget:     highlightPuuid != "" && who.Puuid == highlightPuuid,
		})
	}
	return detail
}

// GetItemIcons resolves item IDs to icon data URIs (cached; failures are skipped).
func (svc *MatchHistoryService) GetItemIcons(ids []int) map[int]string {
	result := make(map[int]string)
	paths := svc.loadItemPaths()
	if paths == nil {
		return result
	}
	cfg, ok := config.Instance()
	if !ok {
		return result
	}

	var (
		wg   sync.WaitGroup
		resM sync.Mutex
	)
	seen := make(map[int]bool, len(ids))
	for _, id := range ids {
		if id <= 0 || seen[id] {
			continue
		}
		seen[id] = true

		svc.cacheMu.RLock()
		cached, hit := svc.itemIcons[id]
		svc.cacheMu.RUnlock()
		if hit {
			result[id] = cached
			continue
		}
		path, known := paths[id]
		if !known {
			continue
		}

		wg.Add(1)
		go func(id int, path string) {
			defer wg.Done()
			data, err := cfg.SendHttpRequest(path, http.MethodGet)
			if err != nil || len(data) == 0 {
				slog.Debug("item icon unavailable", "id", id, "err", err)
				return
			}
			uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
			svc.cacheMu.Lock()
			svc.itemIcons[id] = uri
			svc.cacheMu.Unlock()
			resM.Lock()
			result[id] = uri
			resM.Unlock()
		}(id, path)
	}
	wg.Wait()
	return result
}

func (svc *MatchHistoryService) loadItemPaths() map[int]string {
	svc.cacheMu.RLock()
	paths := svc.itemPaths
	svc.cacheMu.RUnlock()
	if paths != nil {
		return paths
	}

	body, err := svc.cachedGet(itemsPath, 24*time.Hour)
	if err != nil {
		slog.Warn("load items metadata failed", "err", err)
		return nil
	}
	var list []struct {
		ID       int    `json:"id"`
		IconPath string `json:"iconPath"`
	}
	if json.Unmarshal(body, &list) != nil {
		return nil
	}
	paths = make(map[int]string, len(list))
	for _, it := range list {
		paths[it.ID] = it.IconPath
	}

	svc.cacheMu.Lock()
	svc.itemPaths = paths
	svc.cacheMu.Unlock()
	return paths
}

// GetCurrentPuuid returns the logged-in summoner's PUUID ("" if unknown).
func (svc *MatchHistoryService) GetCurrentPuuid() string {
	body, err := svc.cachedGet("/lol-summoner/v1/current-summoner", slowCacheTTL)
	if err != nil {
		return ""
	}
	var s struct {
		Puuid string `json:"puuid"`
	}
	if json.Unmarshal(body, &s) != nil {
		return ""
	}
	return s.Puuid
}
