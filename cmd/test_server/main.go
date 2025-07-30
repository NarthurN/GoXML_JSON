package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// JSONUser - структура пользователя в JSON формате (как в вашем задании)
type JSONUser struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	AgeGroup string `json:"age_group"`
}

// ResponseData - структура ответа сервера
type ResponseData struct {
	Status      string     `json:"status"`
	Message     string     `json:"message"`
	ReceivedAt  string     `json:"received_at"`
	ProcessedBy string     `json:"processed_by"`
	UserCount   int        `json:"user_count"`
	Users       []JSONUser `json:"users"`
}

func main() {
	// Настройка роутов
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health", handleHealth)

	fmt.Println("🚀 Тестовый сервер для получения JSON пользователей")
	fmt.Println("📍 Адрес: http://localhost:8081")
	fmt.Println("🎯 Endpoint для тестирования: http://localhost:8081/users")
	fmt.Println("💡 Этот сервер принимает POST запросы с JSON массивом пользователей")
	fmt.Println("📋 Ожидаемый формат: [{\"id\":\"1\",\"full_name\":\"Иван Иванов\",\"email\":\"ivan@example.com\",\"age_group\":\"от 25 до 35\"}]")
	fmt.Println("🔄 Сервер логирует все запросы и возвращает подтверждение")
	fmt.Println("")
	fmt.Println("🌐 Запуск сервера на порту 8081...")

	// Запуск HTTP сервера
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// handleUsers - основной обработчик для endpoint /users
func handleUsers(w http.ResponseWriter, r *http.Request) {
	// Логируем входящий запрос
	fmt.Printf("\n🔥 [%s] Получен %s запрос на /users от %s\n",
		time.Now().Format("15:04:05"), r.Method, r.RemoteAddr)

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		fmt.Printf("❌ Неподдерживаемый метод: %s (ожидается POST)\n", r.Method)
		http.Error(w, "Метод не поддерживается. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	// Логируем заголовки запроса
	fmt.Println("📋 Заголовки запроса:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("   %s: %s\n", name, value)
		}
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("❌ Ошибка при чтении тела запроса: %v\n", err)
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		fmt.Println("❌ Тело запроса пустое")
		http.Error(w, "Тело запроса не может быть пустым", http.StatusBadRequest)
		return
	}

	fmt.Printf("📄 Получено тело запроса (%d байт):\n%s\n", len(body), string(body))

	// Парсим JSON с пользователями
	var users []JSONUser
	if err := json.Unmarshal(body, &users); err != nil {
		fmt.Printf("❌ Ошибка при парсинге JSON: %v\n", err)
		fmt.Printf("📄 Некорректный JSON: %s\n", string(body))
		http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Printf("✅ JSON успешно распарсен. Получено %d пользователей:\n", len(users))

	// Выводим информацию о каждом пользователе
	for i, user := range users {
		fmt.Printf("   👤 Пользователь %d:\n", i+1)
		fmt.Printf("      • ID: %s\n", user.ID)
		fmt.Printf("      • Полное имя: %s\n", user.FullName)
		fmt.Printf("      • Email: %s\n", user.Email)
		fmt.Printf("      • Возрастная группа: %s\n", user.AgeGroup)
		fmt.Println()
	}

	// Создаем ответ
	response := ResponseData{
		Status:      "success",
		Message:     fmt.Sprintf("Успешно получено и обработано %d пользователей", len(users)),
		ReceivedAt:  time.Now().Format("2006-01-02 15:04:05"),
		ProcessedBy: "Test Server 8081",
		UserCount:   len(users),
		Users:       users,
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Processed-By", "Test-Server-8081")
	w.Header().Set("X-Processing-Time", time.Now().Format(time.RFC3339))
	w.WriteHeader(http.StatusOK)
	// Отправляем JSON ответ
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("❌ Ошибка при кодировании ответа: %v\n", err)
		return
	}

	fmt.Printf("✅ Ответ успешно отправлен клиенту\n")
	fmt.Printf("🎉 Обработка запроса завершена успешно!\n")
	fmt.Println(strings.Repeat("=", 60))
}

// handleRoot - обработчик корневого пути для информации
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n📍 [%s] Запрос на корневой путь от %s\n",
		time.Now().Format("15:04:05"), r.RemoteAddr)

	info := `
🚀 Тестовый сервер для endpoint /users

📋 Информация:
   • Порт: 8081
   • Поддерживаемые методы: POST
   • Endpoint: /users
   • Формат данных: JSON

🎯 Использование:
   POST http://localhost:8081/users
   Content-Type: application/json

   Пример тела запроса:
   [
     {
       "id": "1",
       "full_name": "Иван Иванов",
       "email": "ivan@example.com",
       "age_group": "от 25 до 35"
     }
   ]

🔗 Доступные endpoints:
   • GET  /        - Эта страница с информацией
   • POST /users   - Обработка JSON пользователей
   • GET  /health  - Проверка состояния сервера

💡 Сервер логирует все запросы в консоль для отладки.
`

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, info)
}

// handleHealth - проверка состояния сервера
func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n💊 [%s] Health check от %s\n",
		time.Now().Format("15:04:05"), r.RemoteAddr)

	health := map[string]interface{}{
		"status":    "ok",
		"service":   "Test Server 8081",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
		"endpoints": map[string]string{
			"/":       "GET  - Информация о сервере",
			"/users":  "POST - Обработка JSON пользователей",
			"/health": "GET  - Проверка состояния",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)

	fmt.Println("✅ Health check выполнен")
}
