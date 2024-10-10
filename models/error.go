package models

type SessionError struct {
	Message     string `json:"message"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	Errors      string `json:"errors"`
}

func (s SessionError) Error() string {
	return s.Message
}

type UserAlreadyExists struct {
	SessionError
}

type UserNotFound struct {
	SessionError
}

func (e UserNotFound) Error() string {
	return e.Message
}

type InvalidPassword struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e InvalidPassword) Error() string {
	return e.Message
}
