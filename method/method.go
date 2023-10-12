package method

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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

var exPath string
var apiBaseURL string
var apiTimeout int
var isbns []string

var wg sync.WaitGroup

func init() {
	// Load Environment Variables
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath = filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiBaseURL = os.Getenv("API_PROTOCOL") + "://" + os.Getenv("API_HOST") + ":" + os.Getenv("API_PORT") + "/api"
	apiTimeout, _ = strconv.Atoi(os.Getenv("API_TIMEOUT"))
}

func CallCodebaseOne() *model.BookResponse {
	resultChannel := make(chan model.BookResponse, 1)
	wg.Add(1)

	go callBookIndex(resultChannel)
	
	wg.Wait()
	result := <- resultChannel

	close(resultChannel)
	return &result
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
	_ = json.Unmarshal(resBody, &bookResponse)

	resultChannel <- bookResponse
}

func ConvertIsbn(bookResponse *model.BookResponse) *model.BookResponse {
	for i := 0; i < len(bookResponse.Data); i++ {
		if len(bookResponse.Data[i].Isbn13) != 0 && len(bookResponse.Data[i].Isbn10) == 0{
			bookResponse.Data[i].Isbn10 = convertIsbn13ToIsbn10(bookResponse.Data[i].Isbn13)
		} else if  len(bookResponse.Data[i].Isbn10) != 0 && len(bookResponse.Data[i].Isbn13) == 0 {
			bookResponse.Data[i].Isbn13 = convertIsbn10ToIsbn13(bookResponse.Data[i].Isbn10)
		}
	}
	return bookResponse
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

	isbns = append(isbns, newIsbn)

	return newIsbn
}

func convertIsbn10ToIsbn13(isbn string) string {
	isbn = "978" + isbn
	var isbnArr = strings.Split(isbn, "")
	var isbnSum int
	var newIsbn string
	for i := 0; i < len(isbnArr)-1; i++ {
		var multiplier int = 3
		if i % 2 == 0 {
			multiplier = 1
		}
		num, _ := strconv.Atoi(isbnArr[i])
		newIsbn += strconv.Itoa(num)
		isbnSum += num * multiplier
	}
	checkDigit := 0
	remainder := isbnSum % 10
	
	if remainder != 0 {
		checkDigit = 10 - remainder 
	}
	newIsbn += strconv.Itoa(checkDigit)

	isbns = append(isbns, newIsbn)

	return newIsbn
}

func UpdateBookData(bookResponse *model.BookResponse) {

	for i := 0; i <  len(bookResponse.Data); i++ {
		wg.Add(1)
		go updateRoutine(bookResponse.Data[i])
	}

	wg.Add(1)
	go appendIsbnToCSV()

	wg.Wait()
}

func updateRoutine(book model.Book){
	defer wg.Done()

	var bookUpdate model.BookUpdate
	bookUpdate.ID = book.ID
	bookUpdate.Title = book.Title
	bookUpdate.Isbn13 = book.Isbn13
	bookUpdate.Isbn10 = book.Isbn10
	bookUpdate.PublicationYear = book.PublicationYear
	bookUpdate.PublisherID = book.PublisherID
	bookUpdate.ImageURL = book.ImageURL
	bookUpdate.Edition = book.Edition
	bookUpdate.ListPrice = book.ListPrice
	authorIDsString := book.AuthorIDs
	authorIDsString = strings.Replace(authorIDsString, "[", "", 1)
	authorIDsString = strings.Replace(authorIDsString, "]", "", 1)

	for _, v := range  strings.Split(authorIDsString, ",") {
		num, _ := strconv.ParseUint(v, 10, 64)
		bookUpdate.AuthorIDs = append(bookUpdate.AuthorIDs, num)
	}

	marshalled, err := json.Marshal(bookUpdate)

	if err != nil {
		log.Fatalf("impossible to marshall: %s", err)
	}

	req, err := http.NewRequest(http.MethodPatch, apiBaseURL + "/books/" + strconv.FormatUint(book.ID, 10), bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(apiTimeout) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)

	defer res.Body.Close()
}

func appendIsbnToCSV(){
	defer wg.Done()

	csvFile, err := os.OpenFile(exPath + "/output/isbn.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer csvFile.Close()

    reader := csv.NewReader(csvFile)
    existingCsv, _ := reader.ReadAll()

    existingCsv = append(existingCsv, isbns)

    csvFile, err = os.OpenFile(exPath + "/output/isbn.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer csvFile.Close()

    w := csv.NewWriter(csvFile)

    if err := w.WriteAll(existingCsv); err != nil {
        panic(err)
    }

    w.Flush()
}