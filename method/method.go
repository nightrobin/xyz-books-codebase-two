package method

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"xyz-books-codebase-two/model"

	"github.com/joho/godotenv"
)

var apiBaseURL string
var apiTimeout int

var wg sync.WaitGroup

func init() {
	// Load Environment Variables
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiBaseURL = os.Getenv("API_PROTOCOL") + "://" + os.Getenv("API_HOST") + ":" + os.Getenv("API_PORT") + "/api"
	apiTimeout, _ = strconv.Atoi(os.Getenv("API_TIMEOUT"))
}

func CallCodebaseOne() model.BookResponse {
	resultChannel := make(chan model.BookResponse, 1)
	wg.Add(1)

	go callBookIndex(resultChannel)
	
	wg.Wait()
	result := <- resultChannel

	close(resultChannel)
	return result
}

func callBookIndex(resultChannel chan model.BookResponse) {
	defer wg.Done()

	req, err := http.NewRequest(http.MethodGet, apiBaseURL + "/books", nil)
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(apiTimeout) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	var bookResponse model.BookResponse
	err = json.Unmarshal(resBody, &bookResponse)
	if err != nil {
		fmt.Println(err.Error()) 
	}
	resultChannel <- bookResponse
}

func ConvertIsbn(bookResponse model.BookResponse) {
	for _, book := range bookResponse.Data {
		if len(book.Isbn13) != 0 && len(book.Isbn10) == 0{
			book.Isbn10 = convertIsbn13ToIsbn10(book.Isbn13)
			fmt.Println(book)
		} else {
			convertIsbn10ToIsbn13()
		}
	}
}

func convertIsbn13ToIsbn10(isbn string) string {
	var isbnArr = strings.Split(isbn, "")
	var isbnSum int
	var newIsbn string
	for i := 3; i < len(isbnArr)-1; i++ {
		var num int
		num, _ = strconv.Atoi(isbnArr[i])
		isbnSum += (13 - i) * num
		newIsbn += strconv.Itoa(num)
	}

	checkDigit := 11 - (isbnSum % (isbnSum / 11)) 
	checkDigitChar := strconv.Itoa(checkDigit)
	if checkDigitChar == "10" {
		checkDigitChar = "X"
	}
	newIsbn += checkDigitChar

	return newIsbn
}

func convertIsbn10ToIsbn13() {

}