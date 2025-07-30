package models

import "encoding/xml"

// XMLUsers - структура для хранения пользователей из XML
type XMLUsers struct {
	XMLName xml.Name  `xml:"users"` // Имя корневого элемента
	Users   []XMLUser `xml:"user"`  // Массив пользователей
}

// XMLUser - структура для хранения одного пользователя из XML
type XMLUser struct {
	ID    string `xml:"id,attr"` // ID пользователя
	Name  string `xml:"name"`    // Имя пользователя
	Email string `xml:"email"`   // Email пользователя
	Age   int    `xml:"age"`     // Возраст пользователя
}

// JSONUser - структура для хранения одного пользователя в формате JSON
type JSONUser struct {
	ID       string `json:"id"`        // ID пользователя
	FullName string `json:"full_name"` // Имя пользователя
	Email    string `json:"email"`     // Email пользователя
	AgeGroup string `json:"age_group"` // Возрастная группа пользователя
}
