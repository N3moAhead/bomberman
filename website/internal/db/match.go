package db

import "github.com/N3moAhead/bomberman/website/internal/models"

func GetMatches(from, to int) ([]models.Match, error) {
	var matches []models.Match
	err := Conn.Model(&models.Match{}).Preload("Bot1").Preload("Bot2").Offset(from).Limit(to - from).Find(&matches).Error
	return matches, err
}
