package model

type Department struct {
	ID          int    `db:"id" json: "id"`
	Name        string `db: "name" json: "name"`
	Description string `db: "description" json: "derscription,omitempty"`
}

func NewDepartment(id int, name, description string) *Department {
	return &Department{id, name, description}
}
