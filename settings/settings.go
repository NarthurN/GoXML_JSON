// Настройки приложения
package settings

import "time"

const (
	// AuthKey - ключ для авторизации
	AuthKey = "1234567890"

	// OutputFile - файл для логов
	OutputFile = "provider.log"

	// ClientTimeout - таймаут для HTTP запросов
	ClientTimeout = 10 * time.Second

	// Настроки сервера
	ServerHost = "localhost"
	ServerPort = "8080"

	// ClientURL - базовый URL для HTTP запросов
	ClientURL = "http://localhost:8081/users"
)
