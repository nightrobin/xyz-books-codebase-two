package main

import (
	// "fmt"
	"xyz-books-codebase-two/model"
	"xyz-books-codebase-two/method"
)


func main() {
	bookResponse := &model.BookResponse{}
	bookResponse = method.CallCodebaseOne()
	// fmt.Println(bookResponse)	
	bookResponse = method.ConvertIsbn(bookResponse)
	method.UpdateBookData(bookResponse)
	// fmt.Println(bookResponse)	
}