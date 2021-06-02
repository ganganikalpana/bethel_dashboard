package dto

type NewUserRequest struct {
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Password  string `json:"password"`

	Contact_No     string `json:"contact_no"`
	Address_No     string `json:"address_no"`
	Address_Line01 string `json:"address_line01"`
	Address_Line02 string `json:"address_line02"`
	Address_City   string `json:"address_city"`
	Country        string `json:"country"`
}
