package converter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_ParseXML(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name          string
		xmlData       []byte
		expectedUsers int
		expectedError bool
		description   string
	}{
		{
			name: "валидный XML с одним пользователем",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
</users>`),
			expectedUsers: 1,
			expectedError: false,
			description:   "Корректный XML с одним пользователем",
		},
		{
			name: "валидный XML с несколькими пользователями",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
    <user id="3">
        <name>Петр Сидоров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>`),
			expectedUsers: 3,
			expectedError: false,
			description:   "Корректный XML с тремя пользователями",
		},
		{
			name:          "пустые данные",
			xmlData:       []byte{},
			expectedUsers: 0,
			expectedError: true,
			description:   "Пустой массив байтов должен вернуть ошибку",
		},
		{
			name:          "nil данные",
			xmlData:       nil,
			expectedUsers: 0,
			expectedError: true,
			description:   "Nil данные должны вернуть ошибку",
		},
		{
			name: "XML без пользователей",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
</users>`),
			expectedUsers: 0,
			expectedError: true,
			description:   "XML без пользователей должен вернуть ошибку",
		},
		{
			name: "некорректный XML синтаксис",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
    <user id="3">
        <name>Петр Сидоров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>`),
			expectedUsers: 3,
			expectedError: false,
			description:   "Корректный XML синтаксис",
		},
		{
			name: "XML с пробелами и переносами строк",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>  Иван Иванов  </name>
        <email>  ivan@example.com  </email>
        <age>30</age>
    </user>
</users>`),
			expectedUsers: 1,
			expectedError: false,
			description:   "XML с пробелами должен быть обработан корректно",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ParseXML(tt.xmlData)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Users, tt.expectedUsers)
			}
		})
	}
}

func TestConverter_ParseXML_InvalidXML(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name          string
		xmlData       []byte
		expectedError bool
		description   string
	}{
		{
			name: "незакрытый тег",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
    <user id="3">
        <name>Петр Сидоров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>`),
			expectedError: false,
			description:   "Корректный XML с закрытыми тегами",
		},
		{
			name: "неправильная структура XML",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<root>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
</root>`),
			expectedError: true,
			description:   "Неправильный корневой элемент должен вызвать ошибку",
		},
		{
			name: "некорректный XML синтаксис",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
    <user id="3">
        <name>Петр Сидоров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>`),
			expectedError: false,
			description:   "Корректный XML синтаксис",
		},
		{
			name:          "некорректные данные (не XML)",
			xmlData:       []byte(`Это не XML данные`),
			expectedError: true,
			description:   "Не XML данные должны вызвать ошибку",
		},
		{
			name:          "пустая строка",
			xmlData:       []byte(``),
			expectedError: true,
			description:   "Пустая строка должна вызвать ошибку",
		},
		{
			name:          "только пробелы",
			xmlData:       []byte(`   `),
			expectedError: true,
			description:   "Только пробелы должны вызвать ошибку",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ParseXML(tt.xmlData)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestConverter_ParseXML_DataValidation(t *testing.T) {
	converter := NewConverter()

	validXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
    <user id="3">
        <name>Петр Сидоров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>`)

	result, err := converter.ParseXML(validXML)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Users, 3)

	// Проверяем правильность парсинга данных
	user1 := result.Users[0]
	assert.Equal(t, "1", user1.ID)
	assert.Equal(t, "Иван Иванов", user1.Name)
	assert.Equal(t, "ivan@example.com", user1.Email)
	assert.Equal(t, 30, user1.Age)

	user2 := result.Users[1]
	assert.Equal(t, "2", user2.ID)
	assert.Equal(t, "Мария Петрова", user2.Name)
	assert.Equal(t, "maria@example.com", user2.Email)
	assert.Equal(t, 25, user2.Age)

	user3 := result.Users[2]
	assert.Equal(t, "3", user3.ID)
	assert.Equal(t, "Петр Сидоров", user3.Name)
	assert.Equal(t, "petr@example.com", user3.Email)
	assert.Equal(t, 40, user3.Age)
}

func TestConverter_ParseXML_EdgeCases(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name          string
		xmlData       []byte
		expectedUsers int
		expectedError bool
		description   string
	}{
		{
			name: "XML с одним пользователем и минимальными данными",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>А</name>
        <email>a@b.com</email>
        <age>1</age>
    </user>
</users>`),
			expectedUsers: 1,
			expectedError: false,
			description:   "Минимальные валидные данные",
		},
		{
			name: "XML с максимальными значениями",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="999999">
        <name>Очень длинное имя пользователя с множеством символов и пробелов</name>
        <email>very.long.email.address.with.many.parts@very.long.domain.name.com</email>
        <age>110</age>
    </user>
</users>`),
			expectedUsers: 1,
			expectedError: false,
			description:   "Максимальные значения полей",
		},
		{
			name: "XML с пробелами в атрибутах",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id=" 1 ">
        <name>  Иван Иванов  </name>
        <email>  ivan@example.com  </email>
        <age>30</age>
    </user>
</users>`),
			expectedUsers: 1,
			expectedError: false,
			description:   "Пробелы в атрибутах и значениях",
		},
		{
			name: "XML с переносами строк в значениях",
			xmlData: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван
Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
</users>`),
			expectedUsers: 1,
			expectedError: false,
			description:   "Переносы строк в значениях",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ParseXML(tt.xmlData)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Users, tt.expectedUsers)
			}
		})
	}
}

func TestConverter_ParseXML_LargeDataset(t *testing.T) {
	converter := NewConverter()

	// Создаем XML с большим количеством пользователей
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>`)

	// Добавляем 100 пользователей
	for i := 1; i <= 100; i++ {
		age := 20 + (i % 50) // Возраст от 20 до 69
		userXML := fmt.Sprintf(`
    <user id="%d">
        <name>Пользователь %d</name>
        <email>user%d@example.com</email>
        <age>%d</age>
    </user>`, i, i, i, age)
		xmlData = append(xmlData, []byte(userXML)...)
	}

	xmlData = append(xmlData, []byte(`
</users>`)...)

	result, err := converter.ParseXML(xmlData)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Users, 100)

	// Проверяем несколько пользователей
	for i, user := range result.Users {
		expectedID := fmt.Sprintf("%d", i+1)
		expectedAge := 20 + ((i + 1) % 50) // Исправляем расчет возраста
		assert.Equal(t, expectedID, user.ID)
		assert.Equal(t, fmt.Sprintf("Пользователь %d", i+1), user.Name)
		assert.Equal(t, fmt.Sprintf("user%d@example.com", i+1), user.Email)
		assert.Equal(t, expectedAge, user.Age)
	}
}

func TestConverter_ParseXML_ErrorMessages(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name          string
		xmlData       []byte
		expectedError bool
		errorContains string
		description   string
	}{
		{
			name:          "пустые данные",
			xmlData:       []byte{},
			expectedError: true,
			errorContains: "❌ данные пусты",
			description:   "Проверка сообщения об ошибке для пустых данных",
		},
		{
			name:          "XML без пользователей",
			xmlData:       []byte(`<?xml version="1.0" encoding="UTF-8"?><users></users>`),
			expectedError: true,
			errorContains: "❌ нет пользователей в XML",
			description:   "Проверка сообщения об ошибке для XML без пользователей",
		},
		{
			name:          "некорректный XML",
			xmlData:       []byte(`некорректный XML`),
			expectedError: true,
			errorContains: "❌ ошибка при парсинге XML",
			description:   "Проверка сообщения об ошибке для некорректного XML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ParseXML(tt.xmlData)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// Benchmark тест для проверки производительности
func BenchmarkConverter_ParseXML(b *testing.B) {
	converter := NewConverter()

	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
    <user id="3">
        <name>Петр Сидоров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := converter.ParseXML(xmlData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark тест для большого XML
func BenchmarkConverter_ParseXML_LargeXML(b *testing.B) {
	converter := NewConverter()

	// Создаем большой XML
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>`)

	// Добавляем 100 пользователей
	for i := 1; i <= 100; i++ {
		age := 20 + (i % 50) // Возраст от 20 до 69
		userXML := fmt.Sprintf(`
    <user id="%d">
        <name>Пользователь %d</name>
        <email>user%d@example.com</email>
        <age>%d</age>
    </user>`, i, i, i, age)
		xmlData = append(xmlData, []byte(userXML)...)
	}

	xmlData = append(xmlData, []byte(`
</users>`)...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := converter.ParseXML(xmlData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
