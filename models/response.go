package models

type SuccessResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message     string `json:"message"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	Errors      string `json:"errors"`
}
