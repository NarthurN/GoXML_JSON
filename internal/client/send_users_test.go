package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NarthurN/GoXML_JSON/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SendUsers(t *testing.T) {
	tests := []struct {
		name           string
		users          []models.JSONUser
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedError  bool
		expectedBody   string
	}{
		{
			name: "успешная отправка пользователей",
			users: []models.JSONUser{
				{
					ID:       "1",
					FullName: "Иван Иванов",
					Email:    "ivan@example.com",
					AgeGroup: "от 25 до 35",
				},
				{
					ID:       "2",
					FullName: "Мария Петрова",
					Email:    "maria@example.com",
					AgeGroup: "до 25",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод и заголовки
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Читаем и проверяем тело запроса
				var receivedUsers []models.JSONUser
				err := json.NewDecoder(r.Body).Decode(&receivedUsers)
				require.NoError(t, err)
				assert.Len(t, receivedUsers, 2)
				assert.Equal(t, "1", receivedUsers[0].ID)
				assert.Equal(t, "Иван Иванов", receivedUsers[0].FullName)

				// Отправляем успешный ответ
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "success", "message": "Users processed"}`))
			},
			expectedError: false,
			expectedBody:  `{"status": "success", "message": "Users processed"}`,
		},
		{
			name:  "пустой массив пользователей",
			users: []models.JSONUser{},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var receivedUsers []models.JSONUser
				err := json.NewDecoder(r.Body).Decode(&receivedUsers)
				require.NoError(t, err)
				assert.Len(t, receivedUsers, 0)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "success", "message": "No users to process"}`))
			},
			expectedError: false,
			expectedBody:  `{"status": "success", "message": "No users to process"}`,
		},
		{
			name: "сервер возвращает ошибку 400",
			users: []models.JSONUser{
				{
					ID:       "1",
					FullName: "Иван Иванов",
					Email:    "ivan@example.com",
					AgeGroup: "от 25 до 35",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "Invalid data format"}`))
			},
			expectedError: true,
		},
		{
			name: "сервер возвращает ошибку 500",
			users: []models.JSONUser{
				{
					ID:       "1",
					FullName: "Иван Иванов",
					Email:    "ivan@example.com",
					AgeGroup: "от 25 до 35",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
			},
			expectedError: true,
		},
		{
			name: "сервер недоступен (таймаут)",
			users: []models.JSONUser{
				{
					ID:       "1",
					FullName: "Иван Иванов",
					Email:    "ivan@example.com",
					AgeGroup: "от 25 до 35",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Имитируем долгий ответ
				time.Sleep(2 * time.Second)
				w.WriteHeader(http.StatusOK)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый сервер
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Создаем клиент с URL тестового сервера
			client := &Client{
				URL: server.URL,
				client: &http.Client{
					Timeout: 1 * time.Second, // Короткий таймаут для тестов
				},
			}

			// Выполняем тест
			ctx := context.Background()
			result, err := client.SendUsers(ctx, tt.users)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedBody, string(result))
			}
		})
	}
}

func TestClient_SendUsers_NetworkError(t *testing.T) {
	// Создаем клиент с несуществующим URL
	client := &Client{
		URL: "http://localhost:99999", // Несуществующий порт
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}

	ctx := context.Background()
	result, err := client.SendUsers(ctx, []models.JSONUser{
		{
			ID:       "1",
			FullName: "Иван Иванов",
			Email:    "ivan@example.com",
			AgeGroup: "от 25 до 35",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ошибка при отправке пользователей на сервер")
}

func TestClient_SendUsers_EmptyResponse(t *testing.T) {
	// Создаем тестовый сервер, который не отправляет тело ответа
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Не отправляем тело ответа
	}))
	defer server.Close()

	client := &Client{
		URL: server.URL,
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}

	ctx := context.Background()
	result, err := client.SendUsers(ctx, []models.JSONUser{
		{
			ID:       "1",
			FullName: "Иван Иванов",
			Email:    "ivan@example.com",
			AgeGroup: "от 25 до 35",
		},
	})

	// Должен быть успешным, даже с пустым телом ответа
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", string(result))
}

func TestClient_SendUsers_ContextCancellation(t *testing.T) {
	// Создаем тестовый сервер с задержкой
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Долгая обработка
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client := &Client{
		URL: server.URL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем горутину для отмены контекста
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	result, err := client.SendUsers(ctx, []models.JSONUser{
		{
			ID:       "1",
			FullName: "Иван Иванов",
			Email:    "ivan@example.com",
			AgeGroup: "от 25 до 35",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Benchmark тест для проверки производительности
func BenchmarkClient_SendUsers(b *testing.B) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client := &Client{
		URL: server.URL,
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}

	users := []models.JSONUser{
		{
			ID:       "1",
			FullName: "Иван Иванов",
			Email:    "ivan@example.com",
			AgeGroup: "от 25 до 35",
		},
		{
			ID:       "2",
			FullName: "Мария Петрова",
			Email:    "maria@example.com",
			AgeGroup: "до 25",
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.SendUsers(ctx, users)
		if err != nil {
			b.Fatal(err)
		}
	}
}
