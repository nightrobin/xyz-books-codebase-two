module xyz-books-codebase-two/main

go 1.21.1

require (
	github.com/joho/godotenv v1.5.1 // indirect
	xyz-books-codebase-two/method v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	xyz-books-codebase-two/model v0.0.0-00010101000000-000000000000 // indirect
)

replace xyz-books-codebase-two/model => ./model

replace xyz-books-codebase-two/method => ./method
