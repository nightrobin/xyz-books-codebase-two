package main

import (
	"xyz-books-codebase-two/model"
	"xyz-books-codebase-two/method"
)


func main() {
	bookResponse := &model.BookResponse{}
	bookResponse = method.CallCodebaseOne()
	bookResponse = method.ConvertIsbn(bookResponse)
	method.UpdateBookData(bookResponse)
}