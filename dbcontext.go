package main

import (
	"database/sql"

	"traction_staff_library/models"

	"github.com/itrepablik/itrlog"
)

type Staff struct {
	Id        int8   `json:"Id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TelNo     string `json:"tel_no"`
	IsAdmin   bool   `json:"is_admin"`
}

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

const (
	// Staff Queries
	getStaffQuery    = "SELECT Id, FirstName, LastName, TelNo, IsAdmin FROM staffinfo WHERE Id=$1"
	getStaffsQuery   = "SELECT Id, FirstName, LastName, TelNo, IsAdmin FROM staffinfo LIMIT $1 OFFSET $2"
	updateStaffQuery = "UPDATE staffinfo SET name=$1, price=$2 WHERE id=$3"
	deleteStaffQuery = "DELETE FROM staffinfo WHERE Id=$1"
	createStaffQuery = "INSERT INTO staffinfo(FirstName, LastName, TelNo, IsAdmin) VALUES ($1, $2, $3, $4) RETURNING Id"
)
const (
	// Book Queries
	getBookQuery          = "SELECT Id, Name, Author, ISBN, IsBorrowed, BorrowedByUserId, DateBorrowed, ReturnDate, IsReturned  FROM Book WHERE Id=$1"
	getBooksQuery         = "SELECT Id, Name, Author, ISBN, IsBorrowed, BorrowedByUserId, DateBorrowed, ReturnDate, IsReturned FROM Book LIMIT $1 OFFSET $2"
	updateBookQuery       = "UPDATE Book SET Name=$1, Author=$2, ISBN=$3, IsBorrowed=$4, BorrowedByUserId=$5, DateBorrowed=$6, ReturnDate=$7, IsReturned=$8"
	deleteBookQuery       = "DELETE FROM Book WHERE Id=$1"
	createBookQuery       = "INSERT INTO Book(Name, Author, ISBN, IsBorrowed, BorrowedByUserId, DateBorrowed, ReturnDate, IsReturned) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id"
	getBorrowedBooksQuery = "SELECT Id, Name, Author, ISBN, IsBorrowed, BorrowedByUserId, DateBorrowed, ReturnDate, IsReturned FROM Book WHERE IsBorrowed=true AND IsReturned=false  LIMIT $1 OFFSET $2"
)

func (staff *Staff) getStaff(db *sql.DB) error {
	return db.QueryRow(getStaffQuery, staff.Id).Scan(&staff.Id, &staff.FirstName, &staff.LastName, &staff.TelNo, &staff.IsAdmin)
}

func (staff *Staff) updateProduct(db *sql.DB) error {
	_, err := db.Exec(updateStaffQuery, staff.FirstName, staff.LastName, staff.TelNo, staff.IsAdmin)
	return err
}

func (staff *Staff) deleteProduct(db *sql.DB) error {
	_, err := db.Exec(deleteStaffQuery, staff.Id)
	return err
}

func (staff *Staff) CreateStaff(db *sql.DB) error {
	err := db.QueryRow(
		createStaffQuery,
		staff.FirstName,
		staff.LastName,
		staff.TelNo,
		staff.IsAdmin,
	).Scan(&staff.Id)
	if err != nil {
		itrlog.Error("Error Creating Staff: %e", err)
		return err
	}
	return nil
}

func getStaffs(db *sql.DB, start, count int) ([]models.Staff, error) {
	rows, err := db.Query(getStaffsQuery, count, start)
	if err != nil {
		itrlog.Error("Error occurred ", err)
		return nil, err
	}
	defer rows.Close()

	staffs := []models.Staff{}

	for rows.Next() {
		var staff models.Staff
		if err := rows.Scan(&staff.Id, &staff.FirstName, &staff.LastName, &staff.TelNo, &staff.IsAdmin); err != nil {
			return nil, err
		}
		staffs = append(staffs, staff)
	}

	return staffs, nil
}

func (book *Book) getBook(db *sql.DB) error {
	return db.QueryRow(getBookQuery, book.Id).Scan(
		&book.Id,
		&book.Name,
		&book.Author,
		&book.ISBN,
		&book.IsBorrowed,
		&book.BorrowedByUserId,
		&book.DateBorrowed,
		&book.ReturnDate,
		&book.IsReturned,
	)
}

func getBorrowedBooks(db *sql.DB, start, count int) ([]models.Book, error) {
	rows, err := db.Query(getBorrowedBooksQuery, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []models.Book{}

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.Id,
			&book.Name,
			&book.Author,
			&book.ISBN,
			&book.IsBorrowed,
			&book.BorrowedByUserId,
			&book.DateBorrowed,
			&book.ReturnDate,
			&book.IsReturned); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func getBooks(db *sql.DB, start, count int) ([]models.Book, error) {
	rows, err := db.Query(getBooksQuery, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []models.Book{}

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.Id,
			&book.Name,
			&book.Author,
			&book.ISBN,
			&book.IsBorrowed,
			&book.BorrowedByUserId,
			&book.DateBorrowed,
			&book.ReturnDate,
			&book.IsReturned); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (book *Book) updateBook(db *sql.DB) error {
	_, err := db.Exec(updateBookQuery,
		book.Name,
		book.Author,
		book.ISBN,
		book.IsBorrowed,
		book.BorrowedByUserId,
		book.DateBorrowed,
		book.ReturnDate,
		book.IsReturned,
		book.Id)
	return err
}

func (book *Book) deleteBook(db *sql.DB) error {
	_, err := db.Exec(deleteBookQuery, book.Id)
	return err
}

func (book *Book) CreateBook(db *sql.DB) error {
	err := db.QueryRow(
		createBookQuery,
		book.Name,
		book.Author,
		book.ISBN,
		book.IsBorrowed,
		book.BorrowedByUserId,
		book.DateBorrowed,
		book.ReturnDate,
		book.IsReturned,
	).Scan(&book.Id)
	if err != nil {
		itrlog.Error("Error Creating Staff: %e", err)
		return err
	}
	return nil
}
