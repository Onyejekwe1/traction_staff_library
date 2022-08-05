package models

type Book struct {
	Id               int8   `json:"Id"`
	Name             string `json:"name"`
	Author           string `json:"author"`
	ISBN             string `json:"isbn"`
	IsBorrowed       bool   `json:"is_borrowed"`
	BorrowedByUserId string `json:"borrowed_by_user_id"`
	DateBorrowed     string `json:"date_borrowed"`
	ReturnDate       string `json:"return_date"`
	IsReturned       bool   `json:"is_returned"`
}
