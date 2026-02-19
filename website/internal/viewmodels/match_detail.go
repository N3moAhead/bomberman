package viewmodels

import (
	"encoding/json"

	"github.com/N3moAhead/bombahead/website/internal/message"
	"github.com/N3moAhead/bombahead/website/internal/models"
)

type MatchDetail struct {
	Match   *models.Match
	History *message.GameHistory
}

func NewMatchDetail(match *models.Match) (*MatchDetail, error) {
	var history message.GameHistory
	err := json.Unmarshal(match.History, &history)
	if err != nil {
		return nil, err
	}

	return &MatchDetail{
		Match:   match,
		History: &history,
	}, nil
}
