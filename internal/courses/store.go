package courses

import (
	"fmt"
	"log"
)

// Структура курса
type Course struct {
	ID    int
	Title string
	Price float64
}

// Пример курсов
var courses = []Course{
	{ID: 1, Title: "Курс по Go", Price: 200},
	{ID: 2, Title: "Курс по Python", Price: 250},
}

// Инициализация (можно добавить загрузку из базы данных или файлов)
func Initialize() {
	log.Println("Инициализация курсов...")
	// Например, загрузка из базы данных или других источников
}

// Получение списка курсов
func GetCoursesList() string {
	var list string
	for _, course := range courses {
		list += course.Title + " - " + fmt.Sprintf("Цена: %.2f", course.Price) + "\n"
	}
	return list
}
