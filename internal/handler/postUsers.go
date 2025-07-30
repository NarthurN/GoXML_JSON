package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

// Users - обработчик для POST запроса на /users
func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	h.logger.Log("🙏 Users: начало обработки запроса")

	// Чтение тела запроса
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Logf("❌ Users: ошибка при чтении тела запроса: %v", err)
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		h.logger.Log("❌ Users: тело запроса пустое")
		http.Error(w, "Тело запроса пустое", http.StatusBadRequest)
		return
	}

	// TODO: убрать, если не тестируется
	h.logger.Logf("✅ Users: тело запроса успешно прочитано, размер: %d байт", len(body))
	// h.logger.Logf("✅ Users: тело запроса успешно прочитано, размер: %d байт", len(body))

	// Парсинг XML
	users, err := h.converter.ParseXML(body)
	if err != nil {
		h.logger.Logf("❌ Users: ошибка при парсинге XML: %v", err)
		http.Error(w, "Ошибка при парсинге XML", http.StatusBadRequest)
		return
	}

	h.logger.Logf("✅ Users: XML успешно пропарсен: %v", users)

	// Асинхронно обрабатываем записи Users из XML в JSON
	jsonUsers, err := h.converter.UsersXMLToJSON(users)
	if len(jsonUsers) == 0 {
		if err != nil {
			h.logger.Logf("❌ Нет валидных пользователей для отправки. Ошибки: %v", err)
		} else {
			h.logger.Log("❌ Входной файл не содержал валидных пользователей.")
		}
		http.Error(w, "Не найдено валидных пользователей в предоставленных данных.", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Logf("⚠️ Часть пользователей не прошла валидацию и была пропущена. Ошибки: %v", err)
	}

	h.logger.Logf("✅ Сконвертировано %d валидных пользователей. Начинаем отправку...", len(jsonUsers))

	// Отправляем JSON пользователей на localhost:8080/users
	h.logger.Log("🙏 Users: Отправляем пользователей на сервер")
	bodyResp, err := h.client.SendUsers(r.Context(), jsonUsers)
	if err != nil {
		h.logger.Logf("❌ Users: ошибка при отправке JSON пользователей: %v", err)
		http.Error(w, "Ошибка при отправке JSON пользователей", http.StatusInternalServerError)
		return
	}
	h.logger.Log("✅ Пользователи успешно отправлены на сервер")

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"data":           json.RawMessage(bodyResp),
		"usersProcessed": len(jsonUsers),
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Logf("❌ Users: ошибка при отправке ответа: %v", err)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}
