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

func (c *Client) SendUsers(ctx context.Context, users []models.JSONUser) error {
	c.logger.Log("🙏 SendUsers: Отправляем пользователей на сервер")

	jsonData, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("❌ SendUsers: ошибка при конвертации пользователей в JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("❌ SendUsers: ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("❌ SendUsers: ошибка при отправке пользователей на сервер: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("❌ SendUsers: ошибка при отправке пользователей: %s - %s", resp.Status, string(bodyBytes))
	}

	c.logger.Log("✅ SendUsers: Пользователи успешно отправлены на сервер")
	return nil
}
