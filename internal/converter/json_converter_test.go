package converter

import (
	"fmt"
	"testing"

	"github.com/NarthurN/GoXML_JSON/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_UsersXMLToJSON(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name           string
		users          *models.XMLUsers
		expectedUsers  int
		expectedErrors bool
		description    string
	}{
		{
			name: "успешная конвертация всех пользователей",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "Иван Иванов",
						Email: "ivan@example.com",
						Age:   30,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   22,
					},
					{
						ID:    "3",
						Name:  "Петр Сидоров",
						Email: "petr@example.com",
						Age:   40,
					},
				},
			},
			expectedUsers:  3,
			expectedErrors: false,
			description:    "Все пользователи валидны и должны быть конвертированы",
		},
		{
			name: "пустой массив пользователей",
			users: &models.XMLUsers{
				Users: []models.XMLUser{},
			},
			expectedUsers:  0,
			expectedErrors: true,
			description:    "Пустой массив должен вернуть ошибку",
		},
		{
			name:           "nil пользователи",
			users:          nil,
			expectedUsers:  0,
			expectedErrors: true,
			description:    "Nil должен вернуть ошибку",
		},
		{
			name: "пользователи с пробелами",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    " 1 ",
						Name:  " Иван Иванов ",
						Email: " ivan@example.com ",
						Age:   30,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: false,
			description:    "Пробелы должны быть удалены",
		},
		{
			name: "пользователи с разными возрастными группами",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "Молодой",
						Email: "young@example.com",
						Age:   20,
					},
					{
						ID:    "2",
						Name:  "Средний",
						Email: "middle@example.com",
						Age:   30,
					},
					{
						ID:    "3",
						Name:  "Старший",
						Email: "old@example.com",
						Age:   50,
					},
				},
			},
			expectedUsers:  3,
			expectedErrors: false,
			description:    "Проверка всех возрастных групп",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.UsersXMLToJSON(tt.users)

			if tt.expectedErrors {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedUsers)
			}
		})
	}
}

func TestConverter_UsersXMLToJSON_ValidationErrors(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name           string
		users          *models.XMLUsers
		expectedUsers  int
		expectedErrors bool
		description    string
	}{
		{
			name: "пользователь с пустым ID",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "",
						Name:  "Иван Иванов",
						Email: "ivan@example.com",
						Age:   30,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   25,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: true,
			description:    "Пользователь с пустым ID должен быть пропущен",
		},
		{
			name: "пользователь с пустым именем",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "",
						Email: "ivan@example.com",
						Age:   30,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   25,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: true,
			description:    "Пользователь с пустым именем должен быть пропущен",
		},
		{
			name: "пользователь с пустым email",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "Иван Иванов",
						Email: "",
						Age:   30,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   25,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: true,
			description:    "Пользователь с пустым email должен быть пропущен",
		},
		{
			name: "пользователь с некорректным возрастом (0)",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "Иван Иванов",
						Email: "ivan@example.com",
						Age:   0,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   25,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: true,
			description:    "Пользователь с возрастом 0 должен быть пропущен",
		},
		{
			name: "пользователь с некорректным возрастом (111)",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "Иван Иванов",
						Email: "ivan@example.com",
						Age:   111,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   25,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: true,
			description:    "Пользователь с возрастом 111 должен быть пропущен",
		},
		{
			name: "пользователь с отрицательным возрастом",
			users: &models.XMLUsers{
				Users: []models.XMLUser{
					{
						ID:    "1",
						Name:  "Иван Иванов",
						Email: "ivan@example.com",
						Age:   -5,
					},
					{
						ID:    "2",
						Name:  "Мария Петрова",
						Email: "maria@example.com",
						Age:   25,
					},
				},
			},
			expectedUsers:  1,
			expectedErrors: true,
			description:    "Пользователь с отрицательным возрастом должен быть пропущен",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.UsersXMLToJSON(tt.users)

			assert.Error(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedUsers)
		})
	}
}

func TestConverter_UsersXMLToJSON_AgeGroups(t *testing.T) {
	converter := NewConverter()

	users := &models.XMLUsers{
		Users: []models.XMLUser{
			{
				ID:    "1",
				Name:  "Молодой",
				Email: "young@example.com",
				Age:   20,
			},
			{
				ID:    "2",
				Name:  "Средний",
				Email: "middle@example.com",
				Age:   30,
			},
			{
				ID:    "3",
				Name:  "Старший",
				Email: "old@example.com",
				Age:   50,
			},
		},
	}

	result, err := converter.UsersXMLToJSON(users)

	require.NoError(t, err)
	require.Len(t, result, 3)

	// Проверяем возрастные группы
	assert.Equal(t, "до 25", result[0].AgeGroup)
	assert.Equal(t, "от 25 до 35", result[1].AgeGroup)
	assert.Equal(t, "старше 35", result[2].AgeGroup)
}

func TestConverter_UsersXMLToJSON_DataTransformation(t *testing.T) {
	converter := NewConverter()

	users := &models.XMLUsers{
		Users: []models.XMLUser{
			{
				ID:    "1",
				Name:  "Иван Иванов",
				Email: "ivan@example.com",
				Age:   30,
			},
		},
	}

	result, err := converter.UsersXMLToJSON(users)

	require.NoError(t, err)
	require.Len(t, result, 1)

	// Проверяем правильность трансформации данных
	jsonUser := result[0]
	assert.Equal(t, "1", jsonUser.ID)
	assert.Equal(t, "Иван Иванов", jsonUser.FullName)
	assert.Equal(t, "ivan@example.com", jsonUser.Email)
	assert.Equal(t, "от 25 до 35", jsonUser.AgeGroup)
}

func TestConverter_UsersXMLToJSON_AllInvalidUsers(t *testing.T) {
	converter := NewConverter()

	users := &models.XMLUsers{
		Users: []models.XMLUser{
			{
				ID:    "",
				Name:  "Иван Иванов",
				Email: "ivan@example.com",
				Age:   30,
			},
			{
				ID:    "2",
				Name:  "",
				Email: "maria@example.com",
				Age:   25,
			},
			{
				ID:    "3",
				Name:  "Петр Сидоров",
				Email: "",
				Age:   40,
			},
		},
	}

	result, err := converter.UsersXMLToJSON(users)

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestConverter_UsersXMLToJSON_MixedValidInvalid(t *testing.T) {
	converter := NewConverter()

	users := &models.XMLUsers{
		Users: []models.XMLUser{
			{
				ID:    "1",
				Name:  "Валидный пользователь",
				Email: "valid@example.com",
				Age:   30,
			},
			{
				ID:    "",
				Name:  "Невалидный пользователь",
				Email: "invalid@example.com",
				Age:   25,
			},
			{
				ID:    "3",
				Name:  "Еще один валидный",
				Email: "another@example.com",
				Age:   40,
			},
		},
	}

	result, err := converter.UsersXMLToJSON(users)

	assert.Error(t, err)
	assert.Len(t, result, 2)

	// Проверяем, что только валидные пользователи попали в результат
	ids := make(map[string]bool)
	for _, user := range result {
		ids[user.ID] = true
	}

	assert.True(t, ids["1"])
	assert.True(t, ids["3"])
	assert.False(t, ids[""]) // Невалидный пользователь не должен быть в результате
}

func TestConverter_UsersXMLToJSON_Concurrency(t *testing.T) {
	converter := NewConverter()

	// Создаем большое количество пользователей для тестирования конкурентности
	users := &models.XMLUsers{
		Users: make([]models.XMLUser, 100),
	}

	for i := 0; i < 100; i++ {
		users.Users[i] = models.XMLUser{
			ID:    fmt.Sprintf("%d", i+1),
			Name:  fmt.Sprintf("Пользователь %d", i+1),
			Email: fmt.Sprintf("user%d@example.com", i+1),
			Age:   20 + (i % 50), // Возраст от 20 до 69
		}
	}

	result, err := converter.UsersXMLToJSON(users)

	require.NoError(t, err)
	require.Len(t, result, 100)

	// Проверяем, что все пользователи были обработаны
	for i, jsonUser := range result {
		expectedID := fmt.Sprintf("%d", i+1)
		assert.Equal(t, expectedID, jsonUser.ID)
		assert.Equal(t, fmt.Sprintf("Пользователь %d", i+1), jsonUser.FullName)
		assert.Equal(t, fmt.Sprintf("user%d@example.com", i+1), jsonUser.Email)
		assert.NotEmpty(t, jsonUser.AgeGroup)
	}
}

// Benchmark тест для проверки производительности
func BenchmarkConverter_UsersXMLToJSON(b *testing.B) {
	converter := NewConverter()

	users := &models.XMLUsers{
		Users: []models.XMLUser{
			{
				ID:    "1",
				Name:  "Иван Иванов",
				Email: "ivan@example.com",
				Age:   30,
			},
			{
				ID:    "2",
				Name:  "Мария Петрова",
				Email: "maria@example.com",
				Age:   25,
			},
			{
				ID:    "3",
				Name:  "Петр Сидоров",
				Email: "petr@example.com",
				Age:   40,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := converter.UsersXMLToJSON(users)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark тест для большого количества пользователей
func BenchmarkConverter_UsersXMLToJSON_LargeDataset(b *testing.B) {
	converter := NewConverter()

	users := &models.XMLUsers{
		Users: make([]models.XMLUser, 1000),
	}

	for i := 0; i < 1000; i++ {
		users.Users[i] = models.XMLUser{
			ID:    fmt.Sprintf("%d", i+1),
			Name:  fmt.Sprintf("Пользователь %d", i+1),
			Email: fmt.Sprintf("user%d@example.com", i+1),
			Age:   20 + (i % 50),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := converter.UsersXMLToJSON(users)
		if err != nil {
			b.Fatal(err)
		}
	}
}
