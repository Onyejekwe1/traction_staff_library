package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/itrepablik/itrlog"
	"net/http"
)

type DBConnect struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *DBConnect) Initialize(host string, port int, user, password, dbname, sslmode string) {
	// Organizing our db connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		itrlog.Error("Error occurred", err)
		panic(err)
	}
	defer db.Close()

	// Trying to connect to our DB
	err = db.Ping()
	if err != nil {
		itrlog.Error("Error occurred", err)
		panic(err)
	}

	fmt.Println("Successfully Connected")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *DBConnect) initializeRoutes() {
	a.Router.HandleFunc("/", Home).Methods("GET")
	a.Router.HandleFunc("/staffs", a.getStaffsHandler).Methods("GET")
	a.Router.HandleFunc("/staff/{id:[0-9]+}", a.getStaffHandler).Methods("GET")
	a.Router.HandleFunc("/staff", a.createStaffHandler).Methods("POST")
	a.Router.HandleFunc("/books", a.getBooksHandler).Methods("GET")
	a.Router.HandleFunc("/book", a.createBookHandler).Methods("POST")
	a.Router.HandleFunc("/book/{id:[0-9]+}", a.updateBook).Methods("PUT")
	a.Router.HandleFunc("/book/{id:[0-9]+}", a.deleteBookHandler).Methods("DELETE")
	http.ListenAndServe(":8081", a.Router)
}
