package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/types"
)

const (
	spellsPath    = "/lol-game-data/assets/v1/summoner-spells.json"
	assetMissTTL  = 5 * time.Minute
	assetMaxBytes = 4 << 20
)

// Only these client-served resource trees may be fetched on behalf of the UI.
var assetPrefixes = []string{"/lol-game-data/", "/fe/", "/lol-static-assets/"}

// AllowedAssetPath reports whether path is a plain absolute resource path
// under one of the whitelisted trees (no scheme, host or traversal).
func AllowedAssetPath(path string) bool {
	if strings.Contains(path, "..") || strings.ContainsAny(path, "\\\r\n") || strings.Contains(path, "://") {
		return false
	}
	for _, p := range assetPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// mimeFor guesses the image type from the path extension.
func mimeFor(path string) string {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(lower, ".webp"):
		return "image/webp"
	default:
		return "image/png"
	}
}

// GetAsset returns the first candidate that the client serves as a data URI,
// or "" when none works. Hits are cached for the session; misses for a few
// minutes so a wrong path is not retried on every render.
func (svc *MatchHistoryService) GetAsset(candidates []string) string {
	key := strings.Join(candidates, "|")

	svc.cacheMu.RLock()
	uri, hit := svc.assetURIs[key]
	missUntil, missed := svc.assetMisses[key]
	svc.cacheMu.RUnlock()
	if hit {
		return uri
	}
	if missed && svc.now().Before(missUntil) {
		return ""
	}

	cfg, ok := config.Instance()
	if !ok {
		return ""
	}

	for _, path := range candidates {
		if !AllowedAssetPath(path) {
			slog.Warn("asset path rejected", "path", path)
			continue
		}
		data, err := cfg.SendHttpRequest(path, http.MethodGet)
		if err != nil || len(data) == 0 || len(data) > assetMaxBytes {
			continue
		}
		uri = "data:" + mimeFor(path) + ";base64," + base64.StdEncoding.EncodeToString(data)

		svc.cacheMu.Lock()
		svc.assetURIs[key] = uri
		svc.cacheMu.Unlock()
		slog.Info("asset resolved", "path", path)
		return uri
	}

	slog.Warn("asset not available from the client", "tried", candidates)
	svc.cacheMu.Lock()
	svc.assetMisses[key] = svc.now().Add(assetMissTTL)
	svc.cacheMu.Unlock()
	return ""
}

// RankEmblemCandidates lists where the client may keep the emblem of a tier.
func RankEmblemCandidates(tierKey string) []string {
	tier := strings.ToLower(tierKey)
	if tier == "" || tier == "na" {
		tier = "unranked"
	}
	return []string{
		fmt.Sprintf("/fe/lol-static-assets/images/ranked-emblem/emblem-%s.png", tier),
		fmt.Sprintf("/fe/lol-static-assets/images/ranked-mini-crests/%s.png", tier),
		fmt.Sprintf("/lol-static-assets/v1/ranked-emblem/emblem-%s.png", tier),
	}
}

// PositionIconCandidates lists where the client may keep a lane icon.
func PositionIconCandidates(position string) []string {
	p := strings.ToLower(position)
	return []string{
		fmt.Sprintf("/fe/lol-static-assets/images/position-selector/positions/icon-position-%s.png", p),
		fmt.Sprintf("/fe/lol-static-assets/images/position-selector/positions/icon-position-%s-blue.png", p),
	}
}

type spellMeta struct {
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
}

// GetSpell resolves a summoner spell to the client's own name and icon.
// Falls back to the built-in Chinese name when the client data is unavailable.
func (svc *MatchHistoryService) GetSpell(id int) types.SpellInfo {
	info := types.SpellInfo{ID: id, Name: SpellName(id)}
	if id <= 0 {
		return info
	}

	body, err := svc.cachedGet(spellsPath, 24*time.Hour)
	if err != nil {
		return info
	}
	var list []struct {
		ID int `json:"id"`
		spellMeta
	}
	if json.Unmarshal(body, &list) != nil {
		return info
	}
	for _, s := range list {
		if s.ID != id {
			continue
		}
		if s.Name != "" {
			info.Name = s.Name
		}
		if s.IconPath != "" {
			info.Icon = svc.GetAsset([]string{s.IconPath})
		}
		break
	}
	return info
}
