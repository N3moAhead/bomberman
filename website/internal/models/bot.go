package models

import "gorm.io/gorm"

type Bot struct {
	gorm.Model
	Name          string
	Description   string
	DockerHubUrl  string
	CreatedWithAi bool
	UserID        uint
	User          User
}
