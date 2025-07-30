package converter

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/NarthurN/GoXML_JSON/internal/models"
)

// ParseXML - функция для парсинга XML данных
func (c *Converter) ParseXML(data []byte) (*models.XMLUsers, error) {
	if len(data) == 0 {
		return nil, models.ErrEmptyData
	}

	var users models.XMLUsers
	if err := xml.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("❌ ошибка при парсинге XML: %w", err)
	}

	if len(users.Users) == 0 {
		return nil, models.ErrEmptyUsers
	}

	if err := validateUsers(&users); err != nil {
		return nil, err
	}

	return &users, nil
}

// validateUsers - функция для валидации пользователей
func validateUsers(users *models.XMLUsers) error {
	if users == nil {
		return models.ErrEmptyUsers
	}

	for i := range users.Users {
		u := &users.Users[i]

		// Очистка полей
		u.ID = strings.TrimSpace(u.ID)
		u.Name = strings.TrimSpace(u.Name)
		u.Email = strings.TrimSpace(u.Email)

		// Валидация обязательных полей
		if u.ID == "" {
			return models.ErrEmptyID
		}
		if u.Name == "" {
			return models.ErrEmptyName
		}
		if u.Email == "" {
			return models.ErrEmptyEmail
		}
		if u.Age <= 0 || u.Age > 110 {
			return models.ErrInvalidAge
		}
	}
	return nil
}
