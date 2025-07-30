package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverter_GetAgeGroup(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		age      int
		expected string
	}{
		{
			name:     "возраст меньше 25 - группа 'до 25'",
			age:      0,
			expected: AgeGroupYang,
		},
		{
			name:     "возраст 1 - группа 'до 25'",
			age:      1,
			expected: AgeGroupYang,
		},
		{
			name:     "возраст 10 - группа 'до 25'",
			age:      10,
			expected: AgeGroupYang,
		},
		{
			name:     "возраст 24 - группа 'до 25'",
			age:      24,
			expected: AgeGroupYang,
		},
		{
			name:     "возраст 25 - группа 'от 25 до 35'",
			age:      25,
			expected: AgeGroupMiddle,
		},
		{
			name:     "возраст 30 - группа 'от 25 до 35'",
			age:      30,
			expected: AgeGroupMiddle,
		},
		{
			name:     "возраст 35 - группа 'от 25 до 35'",
			age:      35,
			expected: AgeGroupMiddle,
		},
		{
			name:     "возраст 36 - группа 'старше 35'",
			age:      36,
			expected: AgeGroupOld,
		},
		{
			name:     "возраст 50 - группа 'старше 35'",
			age:      50,
			expected: AgeGroupOld,
		},
		{
			name:     "возраст 100 - группа 'старше 35'",
			age:      100,
			expected: AgeGroupOld,
		},
		{
			name:     "отрицательный возраст - группа 'до 25'",
			age:      -10,
			expected: AgeGroupYang,
		},
		{
			name:     "очень большой возраст - группа 'старше 35'",
			age:      1000,
			expected: AgeGroupOld,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.GetAgeGroup(tt.age)
			assert.Equal(t, tt.expected, result, "Возраст: %d, ожидалось: %s, получено: %s", tt.age, tt.expected, result)
		})
	}
}

func TestConverter_GetAgeGroup_BoundaryValues(t *testing.T) {
	converter := NewConverter()

	// Тестируем граничные значения
	boundaryTests := []struct {
		age      int
		expected string
		desc     string
	}{
		{24, AgeGroupYang, "24 года - последний год в группе 'до 25'"},
		{25, AgeGroupMiddle, "25 лет - первый год в группе 'от 25 до 35'"},
		{35, AgeGroupMiddle, "35 лет - последний год в группе 'от 25 до 35'"},
		{36, AgeGroupOld, "36 лет - первый год в группе 'старше 35'"},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.desc, func(t *testing.T) {
			result := converter.GetAgeGroup(tt.age)
			assert.Equal(t, tt.expected, result, "Возраст: %d, ожидалось: %s, получено: %s", tt.age, tt.expected, result)
		})
	}
}

func TestConverter_GetAgeGroup_Constants(t *testing.T) {
	// Проверяем, что константы определены корректно
	assert.Equal(t, "до 25", AgeGroupYang)
	assert.Equal(t, "от 25 до 35", AgeGroupMiddle)
	assert.Equal(t, "старше 35", AgeGroupOld)
}

func TestConverter_GetAgeGroup_Consistency(t *testing.T) {
	converter := NewConverter()

	// Проверяем консистентность результатов для одинаковых входных данных
	testAges := []int{0, 24, 25, 30, 35, 36, 100}

	for _, age := range testAges {
		result1 := converter.GetAgeGroup(age)
		result2 := converter.GetAgeGroup(age)
		result3 := converter.GetAgeGroup(age)

		assert.Equal(t, result1, result2, "Результаты должны быть одинаковыми для возраста %d", age)
		assert.Equal(t, result2, result3, "Результаты должны быть одинаковыми для возраста %d", age)
		assert.Equal(t, result1, result3, "Результаты должны быть одинаковыми для возраста %d", age)
	}
}

func TestConverter_GetAgeGroup_AllGroups(t *testing.T) {
	converter := NewConverter()

	// Проверяем, что все группы возвращаются
	youngAges := []int{0, 1, 10, 15, 20, 24}
	middleAges := []int{25, 26, 30, 34, 35}
	oldAges := []int{36, 40, 50, 60, 70, 80, 90, 100}

	// Проверяем группу "до 25"
	for _, age := range youngAges {
		result := converter.GetAgeGroup(age)
		assert.Equal(t, AgeGroupYang, result, "Возраст %d должен быть в группе 'до 25'", age)
	}

	// Проверяем группу "от 25 до 35"
	for _, age := range middleAges {
		result := converter.GetAgeGroup(age)
		assert.Equal(t, AgeGroupMiddle, result, "Возраст %d должен быть в группе 'от 25 до 35'", age)
	}

	// Проверяем группу "старше 35"
	for _, age := range oldAges {
		result := converter.GetAgeGroup(age)
		assert.Equal(t, AgeGroupOld, result, "Возраст %d должен быть в группе 'старше 35'", age)
	}
}

// Benchmark тест для проверки производительности
func BenchmarkConverter_GetAgeGroup(b *testing.B) {
	converter := NewConverter()
	testAges := []int{0, 24, 25, 30, 35, 36, 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		age := testAges[i%len(testAges)]
		converter.GetAgeGroup(age)
	}
}

// Benchmark тест для каждой возрастной группы отдельно
func BenchmarkConverter_GetAgeGroup_Young(b *testing.B) {
	converter := NewConverter()
	for i := 0; i < b.N; i++ {
		converter.GetAgeGroup(20)
	}
}

func BenchmarkConverter_GetAgeGroup_Middle(b *testing.B) {
	converter := NewConverter()
	for i := 0; i < b.N; i++ {
		converter.GetAgeGroup(30)
	}
}

func BenchmarkConverter_GetAgeGroup_Old(b *testing.B) {
	converter := NewConverter()
	for i := 0; i < b.N; i++ {
		converter.GetAgeGroup(50)
	}
}
