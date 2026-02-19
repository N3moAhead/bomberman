package db

import (
	"github.com/N3moAhead/bombahead/website/internal/models"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := Conn.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	return Conn.Create(user).Error
}

func GetOrCreateUser(user *models.User) (*models.User, error) {
	existingUser, err := GetUserByUsername(user.Username)
	if err == nil {
		if existingUser.AvatarURL != user.AvatarURL {
			existingUser.AvatarURL = user.AvatarURL
			if err := Conn.Save(existingUser).Error; err != nil {
				return nil, err
			}
		}
		return existingUser, nil
	}

	if err := CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetBotsForUser(user *models.User) ([]models.Bot, error) {
	var bots []models.Bot
	err := Conn.Model(&models.Bot{}).Where("user_id = ?", user.ID).Find(&bots).Error
	return bots, err
}
