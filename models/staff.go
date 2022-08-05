package models

type Staff struct {
	Id        int8   `json:"Id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TelNo     string `json:"tel_no"`
	IsAdmin   bool   `json:"is_admin"`
}
