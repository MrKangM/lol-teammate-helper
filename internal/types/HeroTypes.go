package types

type HeroInfo struct {
	Name               string `json:"name"`
	SquarePortraitPath string `json:"squarePortraitPath"`
	IconDataURI        string `json:"iconDataURI,omitempty"`
}
