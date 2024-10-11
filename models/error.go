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
