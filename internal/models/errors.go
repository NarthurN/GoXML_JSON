package models

import "errors"

// Ошибки
var (
	// Ошибки валидации
	ErrEmptyID    = errors.New("❌ пустой id")
	ErrEmptyName  = errors.New("❌ пустое имя")
	ErrEmptyEmail = errors.New("❌ пустой email")
	ErrInvalidAge = errors.New("❌ некорректный возраст")

	// Ошибки парсинга
	ErrEmptyUsers = errors.New("❌ нет пользователей в XML")

	// Ошибки преобразования
	ErrEmptyData  = errors.New("❌ данные пусты")
	ErrNoUsers    = errors.New("❌ нет пользователей")
)
