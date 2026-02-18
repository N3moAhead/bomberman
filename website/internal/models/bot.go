package models

import (
	"github.com/intinig/go-openskill/types"
	"gorm.io/gorm"
)

type Bot struct {
	gorm.Model
	Name          string
	Description   string
	DockerHubUrl  string
	CreatedWithAi bool
	UserID        uint
	User          User
	Mu            float64 `gorm:"default:25.0"`
	Sigma         float64 `gorm:"default:8.333"`
	Score         float64 `gorm:"index"`
	Wins          int64   `gorm:"->;column:wins;-:migration" json:"wins"`
	Losses        int64   `gorm:"->;column:losses;-:migration" json:"losses"`
	Draws         int64   `gorm:"->;column:draws;-:migration" json:"draws"`
	WinRate       float64 `gorm:"-:all" json:"win_rate"`
}

func (b *Bot) BeforeSave(tx *gorm.DB) (err error) {
	// Conservative Rating(Ordinal): Mu - 3 * Sigma
	// OpenSkill
	b.Score = b.Mu - (3.0 * b.Sigma)
	return
}

// Converts rating to OpenSkill Rating
func (b *Bot) ToRating() types.Rating {
	return types.Rating{
		Mu:    b.Mu,
		Sigma: b.Sigma,
	}
}

// Applies openskill rating to bot
func (b *Bot) ApplyRating(r types.Rating) {
	b.Mu = r.Mu
	b.Sigma = r.Sigma
}

func (b *Bot) CalculateWinRate() {
	totalGames := b.Wins + b.Losses + b.Draws
	if totalGames > 0 {
		b.WinRate = (float64(b.Wins) / float64(totalGames)) * 100
	} else {
		b.WinRate = 0
	}
}
