package service

import (
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/niluwats/bethel_dashboard/dbhandler"
	"github.com/niluwats/bethel_dashboard/domain"
	"github.com/niluwats/bethel_dashboard/dto"
	"github.com/niluwats/bethel_dashboard/errs"
)

type AuthService interface {
	Register(dto.NewUserRequest) (*domain.User, *errs.AppError)
	Login(dto.NewLoginRequest) (*dto.LoginResponse, *errs.AppError)
	Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
	VerifyEmail(request dto.EmailVerifyCode) *errs.AppError
	VerifyMobile(request dto.MobileVerifyCode) *errs.AppError
	Recover(request dto.RecoverRequest) *errs.AppError
	ResetPassword(request dto.PasswordReset, evpw, url_email string) *errs.AppError
}
type DefaultAuthService struct {
	repo dbhandler.AuthRepository
}

func (s DefaultAuthService) ResetPassword(req dto.PasswordReset, evpw, url_email string) *errs.AppError {
	err := s.repo.ResetPassword(evpw, url_email, req.Email, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultAuthService) Recover(req dto.RecoverRequest) *errs.AppError {
	err := s.repo.RecoverEmail(req.Email)
	if err != nil {
		return err
	}
	return nil
}
func (s DefaultAuthService) Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {

			var appErr *errs.AppError
			if appErr = s.repo.RefreshTokenExists(request.RefreshToken); appErr != nil {
				return nil, appErr
			}

			var accessToken string
			if accessToken, appErr = domain.NewAccessTokenFromRefreshToken(request.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errs.NewAuthenticationError("invalid token")
	}
	return nil, errs.NewAuthenticationError("cannot generate a new access token until the current one expires")
}
func (s DefaultAuthService) VerifyEmail(req dto.EmailVerifyCode) *errs.AppError {
	verify := domain.VerifyEmail{
		Email: req.Email,
		Code:  strconv.Itoa(req.Code),
	}
	err := s.repo.VerifyEmail(verify)
	if err != nil {
		return err
	}
	return nil
}
func (s DefaultAuthService) VerifyMobile(req dto.MobileVerifyCode) *errs.AppError {
	verify := domain.VerifyMobile{
		Contact_No: req.Mobile,
		Code:       strconv.Itoa(req.Code),
	}
	err := s.repo.VerifyMobileNo(verify)
	if err != nil {
		return err
	}
	return nil
}
func (s DefaultAuthService) Register(req dto.NewUserRequest) (*domain.User, *errs.AppError) {
	if req.Email == "" || req.Password == "" || req.Country == "" || req.Firstname == "" || req.Lastname == "" {
		return nil, errs.NewUnexpectedError("enter all required fields")
	}
	user := domain.User{
		Email:           req.Email,
		Password:        req.Password,
		Email_Verified:  false,
		Mobile_verified: false,
		Activated:       false,
		Prof: domain.Profile{
			Firstame:       req.Firstname,
			Lastname:       req.Lastname,
			Contact_No:     req.Contact_No,
			Address_No:     req.Address_No,
			Address_Line01: req.Address_Line01,
			Address_Line02: req.Address_Line02,
			Address_City:   req.Address_City,
			Country:        req.Country,
		},
		Role: "user",
	}
	response, err := s.repo.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s DefaultAuthService) Login(req dto.NewLoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.Login
	if req.Email == "" || req.Password == "" {
		return nil, errs.NewUnexpectedError("enter all credentials")
	}
	if login, appErr = s.repo.Login(req.Email, req.Password); appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); appErr != nil {
		return nil, appErr
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
func NewAuthService(repo dbhandler.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repo}
}
