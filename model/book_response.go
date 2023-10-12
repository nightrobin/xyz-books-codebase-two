package model

type BookResponse struct {
	Message	string	`json:"message"`
	Count	int64	`json:"count"`
	Page	int64	`json:"page"`
	Data	[]Book	`json:"data"`
	Errors	[]ApiError	`json:"errors"`	
}
