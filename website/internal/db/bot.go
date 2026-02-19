package db

import (
	"github.com/N3moAhead/bombahead/website/internal/models"
)

const statsSelectQuery = `
	bots.*,
	(
		SELECT COUNT(1)
		FROM matches m
		WHERE m.status = 'finished'
		AND m.deleted_at IS NULL
		AND (
			(m.bot1_id = bots.id AND m.winner_state = 'bot1win') OR
			(m.bot2_id = bots.id AND m.winner_state = 'bot2win')
		)
	) as wins,
	(
		SELECT COUNT(1)
		FROM matches m
		WHERE m.status = 'finished'
		AND m.deleted_at IS NULL
		AND (
			(m.bot1_id = bots.id AND m.winner_state = 'bot2win') OR
			(m.bot2_id = bots.id AND m.winner_state = 'bot1win')
		)
	) as losses,
	(
		SELECT COUNT(1)
		FROM matches m
		WHERE m.status = 'finished'
		AND m.deleted_at IS NULL
		AND m.winner_state = 'draw'
		AND (m.bot1_id = bots.id OR m.bot2_id = bots.id)
	) as draws
`

func CreateBot(bot *models.Bot) error {
	return Conn.Create(bot).Error
}

func GetBotByID(id uint) (*models.Bot, error) {
	var bot models.Bot

	err := Conn.
		Select(statsSelectQuery).
		Preload("User").
		Where("bots.id = ?", id).
		First(&bot).Error

	if err != nil {
		return nil, err
	}

	bot.CalculateWinRate()

	return &bot, nil
}

func GetLeaderboard(page int, pageSize int) ([]models.Bot, error) {
	var bots []models.Bot

	offset := (page - 1) * pageSize

	err := Conn.
		Select(statsSelectQuery).
		Order("score desc").
		Limit(pageSize).
		Offset(offset).
		Preload("User").
		Find(&bots).Error

	if err != nil {
		return nil, err
	}

	for i := range bots {
		bots[i].CalculateWinRate()
	}

	return bots, err
}
