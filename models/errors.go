package models

import "strings"

var (
	ErrNotFound          modelError = "models: resource not found"
	ErrIDInvalid         modelError = "models: ID provided was invalid"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailInvalid      modelError = "models: email provided was invalid"
	ErrEmailIsTaken      modelError = "models: email address has already taken"
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrPasswordTooShort  modelError = "models: password must be at least 8 charaters"
	ErrRememberRequired  modelError = "models: remember is required"
	ErrRememberTooShort  modelError = "models: remember token must be at least 32 bytes"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
