package models

import "gorm.io/gorm"

type MatchStatus string

const (
	PENDING  MatchStatus = "pending"
	RUNNING  MatchStatus = "running"
	FINISHED MatchStatus = "finished"
)

type Match struct {
	gorm.Model
	MatchID string
	Bot1ID  uint
	Bot1    Bot
	Bot2ID  uint
	Bot2    Bot
	Status  MatchStatus
}
