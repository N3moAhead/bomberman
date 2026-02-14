package db

import (
	"github.com/N3moAhead/bomberman/website/internal/models"
)

func CreateBot(bot *models.Bot) error {
	return Conn.Create(bot).Error
}

func GetLeaderboard(page int, pageSize int) ([]models.Bot, error) {
	var bots []models.Bot

	offset := (page - 1) * pageSize

	err := Conn.Order("score desc").
		Limit(pageSize).
		Offset(offset).
		Preload("User").
		Find(&bots).Error

	return bots, err
}
