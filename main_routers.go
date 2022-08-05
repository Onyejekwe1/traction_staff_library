package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/itrepablik/itrlog"
	"html/template"
	"net/http"
	"strconv"
	"traction_staff_library/config"
)

//func MainRouters(r *mux.Router) {
//	r.HandleFunc("/", Home).Methods("GET")
//}

type contextData map[string]interface{}

func Home(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles(config.SiteRootTemplate+"/frontend/index.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "Welcome to Maharlikans Code Tutorial Series",
		"PageMetaDesc": config.SiteSlogan,
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
		"Settings":     config.SiteSettings,
	}
	template.Execute(w, data)
}

// send a payload of JSON content
func (a *DBConnect) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// send a JSON error message
func (a *DBConnect) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, map[string]string{"error": message})

	itrlog.Error("App error: code %d, message %s", code, message)
}

func (a *DBConnect) createStaffHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		itrlog.Error("Error Parsing form: %e", err)
		a.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	staff := Staff{}
	staff.FirstName = r.Form.Get("FirstName")
	staff.LastName = r.Form.Get("LastName")
	staff.TelNo = r.Form.Get("TelNo")

	if err := staff.CreateStaff(a.DB); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusCreated, staff)
}

func (a *DBConnect) getStaffHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		itrlog.Error(msg)
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	staffId, _ := strconv.Atoi(r.FormValue("Id"))
	staff := Staff{Id: int8(staffId)}
	if err := staff.getStaff(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := fmt.Sprintf("Product not found. Error: %s", err.Error())
			a.respondWithError(w, http.StatusNotFound, msg)
		default:
			a.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	a.respondWithJSON(w, http.StatusOK, staff)
}

func (a *DBConnect) getStaffsHandler(w http.ResponseWriter, r *http.Request) {
	//count, _ := strconv.Atoi(r.FormValue("count"))
	//start, _ := strconv.Atoi(r.FormValue("start"))
	//
	//if count > 10 || count < 1 {
	//	count = 10
	//}
	//if start < 0 {
	//	start = 0
	//}

	staffs, err := getStaffs(a.DB, 0, 10)
	if err != nil {
		itrlog.Warn("Error occurred ", err)
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusOK, staffs)
}

func (a *DBConnect) createBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		itrlog.Error("Error Parsing form: %e", err)
		a.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	book := Book{}
	book.Name = r.Form.Get("Name")
	book.ISBN = r.Form.Get("ISBN")
	book.Author = r.Form.Get("Author")
	book.IsBorrowed = false

	if err := book.CreateBook(a.DB); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusCreated, book)
}

func (a *DBConnect) getBooksHandler(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	books, err := getBooks(a.DB, start, count)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusOK, books)
}

func (a *DBConnect) getBorrowedBooksHandler(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	books, err := getBorrowedBooks(a.DB, start, count)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusOK, books)
}

func (a *DBConnect) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		itrlog.Error(msg)
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	bookId, _ := strconv.Atoi(r.FormValue("Id"))
	book := Book{Id: int8(bookId)}
	if err := book.deleteBook(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := fmt.Sprintf("Product not found. Error: %s", err.Error())
			a.respondWithError(w, http.StatusNotFound, msg)
		default:
			a.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	a.respondWithJSON(w, http.StatusOK, book)
}

func (a *DBConnect) updateBook(w http.ResponseWriter, r *http.Request) {
	var b Book

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	b.Id = int8(id)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&b); err != nil {
		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	defer r.Body.Close()

	if err := b.updateBook(a.DB); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusOK, b)
}
