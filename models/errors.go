package models

import (
	"strings"
)

var (
	//ErrNotFound is returned when a resource cannot be found
	//in the database
	ErrNotFound modelError = "models: resource not found"

	//ErrPasswordIncorrect is returned when an invalid password
	//is used when attempting to authenticate a user
	ErrPasswordIncorrect modelError = "models: invalid password provided"

	//ErrEmailRequired is returned when email address is not provided
	//when create a user
	ErrEmailRequired modelError = "models: email address is required"

	//ErrEmailInvalid is returned when an invalid format email address
	//is provided when create a user
	ErrEmailInvalid modelError = "models: email address is not valid"

	//ErrEmailTaken is returned when an email address provided was taken
	//by another user on update and create a user
	ErrEmailTaken modelError = "models: email address is already taken"

	//ErrPasswordTooShort is returned when an update and create attempted
	//with a user password that is less than 8 characters.
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

	//ErrPasswordRequired is returned when an update and create attempted
	//with a user password that is empty
	ErrPasswordRequired modelError = "models: password is required"

	//ErrGalleryTitleRequired return when Title of Gallery is empty
	ErrGalleryTitleRequired modelError = "models: Title is required"

	//ErrPasswordHashRequired is return when an update and create without
	//password hash
	ErrPasswordHashRequired modelError = "models: password hash is required"

	//ErrIDInvalid is returned when a invalid ID is provided
	//to a method like Delete
	ErrIDInvalid privateError = "models: ID provided was invalid"

	//ErrRememberTooShort is return when Remember token string convert to len of bytes
	//at least 32
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"

	//ErrRememberHashRequired is return when Remember Hash is empty
	ErrRememberHashRequired privateError = "models: remember hash is required"

	//ErrorUserIDRequired return when UserID of Gallery is zero
	ErrUserIDRequired privateError =  "models: UserID is required"

)

type modelError string
func (e modelError) Error() string{
	return string(e)
}
func (e modelError) Public() string{
	s:= strings.Replace(string(e),"models: ","", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])

	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string{
	return string(e)
}


