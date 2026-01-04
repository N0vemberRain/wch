package model

type SearchFilter struct {
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	Surname      string `db:"surname" json:"surname"`
	DepartmentID int
}
