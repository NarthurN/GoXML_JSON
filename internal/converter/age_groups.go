package converter

const (
	AgeGroupYang = "до 25"
	AgeGroupMiddle = "от 25 до 35"
	AgeGroupOld = "старше 35"
)

// GetAgeGroup - функция для определения возрастной группы
func (c *Converter) GetAgeGroup(age int) string {
	switch {
	case age < 25:
		return AgeGroupYang
	case age >= 25 && age <= 35:
		return AgeGroupMiddle
	default:
		return AgeGroupOld
	}
}
