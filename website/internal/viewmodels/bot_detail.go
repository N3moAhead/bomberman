package viewmodels

import (
	"github.com/N3moAhead/bombahead/website/internal/models"
)

type BotDetail struct {
	Bot     *models.Bot
	Matches []models.Match
}

func NewBotDetail(bot *models.Bot, matches []models.Match) (*BotDetail, error) {
	return &BotDetail{
		Bot:     bot,
		Matches: matches,
	}, nil
}
