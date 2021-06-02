package domain

import "time"

type PwResetParams struct {
	Code    string    `bson:"code"`
	Timeout time.Time `bson:"timeout"`
}
