package converter

import "github.com/NarthurN/GoXML_JSON/internal/models"

func (c *Converter) UserXMLToJSON(user models.XMLUser) models.JSONUser {
	return models.JSONUser{
		ID:       user.ID,
		FullName: user.Name,
		Email:    user.Email,
		AgeGroup: c.GetAgeGroup(user.Age),
	}
}

func (c *Converter) UsersXMLToJSON(users models.XMLUsers) []models.JSONUser {
	jsonUsers := make([]models.JSONUser, len(users.Users))
	for i, user := range users.Users {
		jsonUsers[i] = c.UserXMLToJSON(user)
	}
	return jsonUsers
}
