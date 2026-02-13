package db

import "github.com/N3moAhead/bomberman/website/internal/models"

func CreateBot(bot *models.Bot) error {
	return Conn.Create(bot).Error
}
