package converter

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/NarthurN/GoXML_JSON/internal/models"
)

// job определяет задание для воркера, включая исходный индекс.
type job struct {
	index int
	user  models.XMLUser
}

// result содержит результат обработки одного задания.
type result struct {
	index    int
	jsonUser models.JSONUser
	err      error
}

func (c *Converter) UserXMLToJSON(user models.XMLUser) models.JSONUser {
	return models.JSONUser{
		ID:       user.ID,
		FullName: user.Name,
		Email:    user.Email,
		AgeGroup: c.GetAgeGroup(user.Age),
	}
}

// UsersXMLToJSON асинхронно конвертирует срез пользователей XML в JSON.
func (c *Converter) UsersXMLToJSON(users *models.XMLUsers) ([]models.JSONUser, error) {
	if users == nil || len(users.Users) == 0 {
		return nil, models.ErrNoUsers
	}

	// Конвертируем XML в JSON асинхронно
	workerCount := runtime.NumCPU()
	if len(users.Users) < workerCount {
		workerCount = len(users.Users)
	}
	jobs := make(chan job, len(users.Users))
	results := make(chan result, len(users.Users))

	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				validatedUser, err := validateUser(job.user)
				if err != nil {
					results <- result{index: job.index, err: fmt.Errorf("user с ID #%d: %w", job.index, err)}
					continue
				}

				jsonUser := models.JSONUser{
					ID:       validatedUser.ID,
					FullName: validatedUser.Name,
					Email:    validatedUser.Email,
					AgeGroup: c.GetAgeGroup(validatedUser.Age),
				}
				results <- result{index: job.index, jsonUser: jsonUser}
			}
		}()
	}

	for i, user := range users.Users {
		jobs <- job{index: i, user: user}
	}
	close(jobs)

	wg.Wait()
	close(results)

	jsonUsers := make([]models.JSONUser, len(users.Users))
	validationErrors := make([]error, 0, len(users.Users))

	for res := range results {
		if res.err != nil {
			validationErrors = append(validationErrors, res.err)
			continue
		}
		jsonUsers[res.index] = res.jsonUser
	}

	// Отсеиваем нулевые структуры, которые остались на месте пользователей с ошибками
	finalUsers := make([]models.JSONUser, 0, len(users.Users)-len(validationErrors))
	for _, u := range jsonUsers {
		if u.ID != "" {
			finalUsers = append(finalUsers, u)
		}
	}

	if len(validationErrors) > 0 {
		return finalUsers, errors.Join(validationErrors...)
	}

	return finalUsers, nil
}

// validateUser - функция для валидации пользователей
func validateUser(user models.XMLUser) (models.XMLUser, error) {
	user = cleanFromSpaces(user)

	if user.ID == "" {
		return user, models.ErrEmptyID
	}
	if user.Name == "" {
		return user, models.ErrEmptyName
	}
	if user.Email == "" {
		return user, models.ErrEmptyEmail
	}
	if user.Age <= 0 || user.Age > 110 {
		return user, models.ErrInvalidAge
	}
	return user, nil
}

// cleanFromSpaces - функция для очистки полей от пробелов
func cleanFromSpaces(user models.XMLUser) models.XMLUser {
	user.ID = strings.TrimSpace(user.ID)
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	return user
}
