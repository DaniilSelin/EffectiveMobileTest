package models

type Person struct {
	ID			 string     `json:"id,omitempty"`
	Name 		 string 	`json:"name,omitempty"`
	Surname 	 string     `json:"surname"`
	Patronymic   string 	`json:"patronymic"`
	Gender	     string 	`json:"gender,omitempty"`
	Age          int 		`json:"age,omitempty"`
	Nationality  string 	`json:"nationality,omitempty"`
}

type PersonFilters struct {
	ID			 string 	`json:"id,omitempty"`
	Name 		 string 	`json:"name,omitempty"`
	Surname 	 string     `json:"surname,omitempty"`
	Patronymic   string 	`json:"patronymic,omitempty"`
	Gender	     string 	`json:"gender,omitempty"`
	Age          int 		`json:"age,omitempty"`
	Nationality  string 	`json:"nationality,omitempty"`
}