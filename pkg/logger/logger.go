// Пакет для логирования
package logger

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/NarthurN/GoXML_JSON/settings"
)

// Logger - потокобезопасный логгер
type Logger struct {
	mu     sync.Mutex
	file   *os.File
	logger *log.Logger
}

// New создает новый логгер, который пишет в OutputFile.
func New() (*Logger, error) {
	file, err := os.OpenFile(settings.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	l := &Logger{
		file:   file,
		logger: log.New(io.MultiWriter(file), "", log.LstdFlags),
	}
	return l, nil
}

// Log пишет строку в лог с временной меткой.
func (l *Logger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Println(msg)
}

// Logf пишет форматированное сообщение в лог.
func (l *Logger) Logf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf(format, args...)
}

// Close закрывает файл логов.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
