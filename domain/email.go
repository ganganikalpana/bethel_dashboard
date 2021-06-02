package domain

type VerifyEmail struct {
	Email string `bson:"email" json:"email"`
	Code  string    `bson:"code" json:"code"`
}
