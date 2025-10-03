package storage

import "github.com/avinendra08/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, gender string, contact int, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}