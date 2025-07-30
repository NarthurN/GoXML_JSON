package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/NarthurN/GoXML_JSON/internal/models"
)

// Users - обработчик для POST запроса на /users
func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.logger.Log("🙏 Users: начало обработки запроса")

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Logf("❌ Users: ошибка при чтении тела запроса: %v", err)
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		h.logger.Log("❌ Users: тело запроса пустое")
		http.Error(w, "Тело запроса пустое", http.StatusBadRequest)
		return
	}

	h.logger.Logf("✅ Users: тело запроса успешно прочитано: %s", string(body))

	// Парсинг XML
	users, err := h.converter.ParseXML(body)
	if err != nil {
		h.logger.Logf("❌ Users: ошибка при парсинге XML: %v", err)
		http.Error(w, "Ошибка при парсинге XML", http.StatusBadRequest)
		return
	}

	h.logger.Logf("✅ Users: XML успешно пропарсен: %v", users)

	// Асинхронно обрабатываем записи Users в JSON
	//jsonUsers := h.converter.UsersXMLToJSON(*users)
}
