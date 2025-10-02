package controller

import (
	"errors"
	"fmt"
	"lol-teammate-helper/internal/config"
	"net/http"
)

// GetPlayerRankMatches fetches the recent ranked match history for the given PUUID.
func GetPlayerRankMatches(puuid string) ([]byte, error) {
	const logPrefix = "[rankController.GetPlayerRankMatches]"

	if puuid == "" {
		err := errors.New("empty puuid provided")
		fmt.Printf("%s %v\n", logPrefix, err)
		return nil, err
	}

	cfg, ok := config.Instance()
	if !ok {
		err := errors.New("riot credentials are not initialised")
		fmt.Printf("%s %v\n", logPrefix, err)
		return nil, err
	}

	endpoint := fmt.Sprintf("/lol-match-history/v1/products/lol/%s/matches", puuid)
	fmt.Println(logPrefix + " requesting " + endpoint)

	body, err := cfg.SendHttpRequest(endpoint, http.MethodGet)
	if err != nil {
		fmt.Printf("%s failed to request matches: %v\n", logPrefix, err)
		return nil, err
	}

	fmt.Printf("%s received %d bytes\n", logPrefix, len(body))
	return body, nil
}
