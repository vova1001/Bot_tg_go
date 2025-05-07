package courses

import (
	"fmt"
	"log"
)

type Course struct {
	ID    int
	Title string
	Price float64
}

var courses = []Course{
	{ID: 1, Title: "Курс по Go", Price: 200},
	{ID: 2, Title: "Курс по Python", Price: 250},
}

func Initialize() {
	log.Println("Инициализация курсов...")

}

func GetCoursesList() string {
	var list string
	for _, course := range courses {
		list += course.Title + " - " + fmt.Sprintf("Цена: %.2f", course.Price) + "\n"
	}
	return list
}
