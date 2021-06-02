package domain

type VerifyMobile struct {
	Contact_No string `bson:"contact_no"`
	Code       string    `bson:"code"`
}
