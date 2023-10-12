package main

import (
	"xyz-books-codebase-two/method"
)


func main() {
	bookResponse := method.CallCodebaseOne()
	method.ConvertIsbn(bookResponse)	
}