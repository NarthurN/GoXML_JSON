package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NarthurN/GoXML_JSON/internal/handler"
	"github.com/NarthurN/GoXML_JSON/pkg/logger"
	appMiddleware "github.com/NarthurN/GoXML_JSON/pkg/middleware"
	"github.com/NarthurN/GoXML_JSON/settings"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logg, err := logger.New()
	if err != nil {
		log.Fatalf("❌ не удалось создать логгер: %v", err)
	}
	defer logg.Close()

	logg.Log("✅ логгер инциализирован")

	handler := handler.NewHandler(logg)
	logg.Log("✅ обработчик инциализирован")

	// Создаем роутер
	r := chi.NewRouter()

	// Используем middleware от chi для надежности
	r.Use(middleware.RequestID)                       // Добавляет уникальный ID каждому запросу
	r.Use(middleware.RealIP)                          // Использует реальный IP клиента
	r.Use(middleware.Logger)                          // Логирует запросы (от chi, не ваш)
	r.Use(middleware.Recoverer)                       // Перехватывает паники и возвращает 500
	r.Use(middleware.Timeout(settings.ClientTimeout)) // Таймаут на весь запрос

	r.Post("/users", handler.Users)

	// Настройка маршрутов
	// Группируем роуты, которые требуют авторизации
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.Auth(logg))
		r.Post("/users", handler.Users)
	})

	// Простой health-check эндпоинт
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	logg.Log("✅ маршруты настроены")

	// Запуск сервера
	srv := &http.Server{
		Addr:    net.JoinHostPort(settings.ServerHost, settings.ServerPort),
		Handler: r,
	}

	go func() {
		logg.Logf("🚀 Сервер слушает на %s", net.JoinHostPort(settings.ServerHost, settings.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Logf("❌ Ошибка запуска сервера: %v", err)
		}
	}()
	logg.Log("✅ сервер запущен")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Log("🚨 Получен сигнал завершения. Начинаем graceful shutdown...")

	// Даем 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), settings.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Logf("❌ Ошибка при graceful shutdown: %v", err)
	}

	logg.Log("✅ Сервер успешно остановлен.")
}
