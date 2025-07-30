package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NarthurN/GoXML_JSON/internal/models"
)

func (c *Client) SendUsers(ctx context.Context, users []models.JSONUser) ([]byte, error) {
	jsonData, err := json.Marshal(users)
	if err != nil {
		return nil, fmt.Errorf("❌ SendUsers: ошибка при конвертации пользователей в JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("❌ SendUsers: ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ SendUsers: ошибка при отправке пользователей на сервер: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("❌ SendUsers: ошибка чтения тела ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("❌ SendUsers: получен неверный статус: %s - %s", resp.Status, string(bodyBytes))
	}

	return bodyBytes, nil
}
