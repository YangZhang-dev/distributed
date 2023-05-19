package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	Id        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var res float32
	for _, grade := range s.Grades {
		res += grade.Score
	}
	return res / float32(len(s.Grades))
}

type Students []Student

var (
	students      Students
	studentsMutex sync.Mutex
)

func (ss Students) GetById(id int) (*Student, error) {
	for i := range ss {
		if ss[i].Id == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("student with id %v not found", id)
}

type GradeType string

const (
	GradeExam = GradeType("Exam")
	GradeQuiz = GradeType("Quiz")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
