package storage

type Storage interface{
	CreateStudent(name string, gender string, contact int,email string,age int)(int64,error)
}