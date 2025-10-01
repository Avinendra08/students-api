package types

type Student struct{
	Id int
	Name string `validate:"required"`
	Gender string `validate:"required"`
	Contact int `validate:"required"`
	Email string `validate:"required"`
	Age int `validate:"required"`
}