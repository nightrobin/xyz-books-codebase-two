package model

type ApiError struct {
	Param   string	`json:"param"`	
	Message string	`json:"message"`	
}