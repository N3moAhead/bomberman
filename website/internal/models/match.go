package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MatchStatus string

const (
	PENDING  MatchStatus = "pending"
	RUNNING  MatchStatus = "running"
	FINISHED MatchStatus = "finished"
)

type WinnerState string

const (
	BOT1WIN WinnerState = "bot1win"
	BOT2WIN WinnerState = "bot2win"
	DRAW    WinnerState = "draw"
)

type Match struct {
	gorm.Model
	MatchID     string
	Bot1ID      uint
	Bot1        Bot
	Bot2ID      uint
	Bot2        Bot
	WinnerState WinnerState
	Status      MatchStatus
	History     datatypes.JSON
}
