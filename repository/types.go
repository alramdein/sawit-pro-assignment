// This file contains types that are used in the repository layer.
package repository

type InsertUserInput struct {
	Id       string
	Name     string
	Phone    string
	Password string
}

type InsertUserOutput struct {
	Id string
}
