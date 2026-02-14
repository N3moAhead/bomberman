package db

import "github.com/N3moAhead/bomberman/website/internal/models"

func GetMatches(from, to int) ([]models.Match, error) {
	var matches []models.Match
	err := Conn.Model(&models.Match{}).Preload("Bot1").Preload("Bot2").Offset(from).Limit(to - from).Find(&matches).Error
	return matches, err
}

func GetMatchByMatchID(matchID string) (*models.Match, error) {
	var match models.Match
	err := Conn.Model(&models.Match{}).Preload("Bot1").Preload("Bot2").Where("match_id = ?", matchID).First(&match).Error
	return &match, err
}
