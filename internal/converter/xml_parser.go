package converter

import (
	"encoding/xml"
	"fmt"

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

	return &users, nil
}
