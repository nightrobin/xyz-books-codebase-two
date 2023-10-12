package model

type BookUpdate struct {
	ID				uint64		`json:"ID"`
	Title			string		`json:"Title"`
	Isbn13			string		`json:"Isbn13"`
	Isbn10			string		`json:"Isbn10"`
	PublicationYear	int16		`json:"PublicationYear"`
	PublisherID		uint64		`json:"PublisherID"`
	ImageURL		string		`json:"ImageURL"`
	Edition			string		`json:"Edition"`
	ListPrice		float32		`json:"ListPrice"`
	AuthorIDs		[]uint64		`json:"AuthorIDs"`
}