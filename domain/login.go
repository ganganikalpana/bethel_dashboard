package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Login struct {
	Email            string `bson:"email"`
	Password         string `bson:"password"`
	Role             string `bson:"role"`
	IsActivated      bool   `bson:"activated"`
	IsEmailVerified  bool   `bson:"email_verified"`
	IsMobileVerified bool   `bson:"mobile_verified"`
	Prof Profile
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	return l.claimsForUser()
}

func (l Login) claimsForUser() AccessTokenClaims {
	return AccessTokenClaims{
		Email: l.Email,
		Role:  l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}
