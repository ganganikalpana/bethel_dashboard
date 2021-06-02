package domain

type User struct {
	Email           string `bson:"email" json:"email"`
	Password        string `bson:"password" json:"password,omitempty"`
	Email_Verified  bool   `bson:"email_verified" json:"email_verified"`
	Mobile_verified bool   `bson:"mobile_verified" json:"mobile_verified"`
	Activated       bool   `bson:"activated" json:"activated"`
	Role            string `bson:"role"`
	Prof            Profile
}
type Profile struct {
	Firstame       string `json:"first_name" bson:"first_name"`
	Lastname       string `json:"last_name" bson:"last_name"`
	Contact_No     string `json:"contact_no" bson:"contact_no"`
	Address_No     string `json:"address_no" bson:"address_no"`
	Address_Line01 string `json:"address_line01" bson:"address_line01"`
	Address_Line02 string `json:"address_line02" bson:"address_line02"`
	Address_City   string `json:"address_city" bson:"address_city"`
	Country        string `json:"country" bson:"country"`
}
