package dbhandler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/niluwats/bethel_dashboard/domain"
	"github.com/niluwats/bethel_dashboard/errs"
	"github.com/niluwats/bethel_dashboard/logger"
	"github.com/niluwats/bethel_dashboard/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AuthRepository interface {
	SaveUser(profile domain.User) (*domain.User, *errs.AppError)
	Login(email, password string) (*domain.Login, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(authToken domain.AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
	VerifyEmail(verify domain.VerifyEmail) *errs.AppError
	VerifyMobileNo(verify domain.VerifyMobile) *errs.AppError
	RecoverEmail(email string) *errs.AppError
	ResetPassword(evpw, url_email, email, newpw, confirmpw string) *errs.AppError
	SaveNode(vm domain.VmAll) (*domain.Organization, *errs.AppError)
	FindIfNodeExists(resgrp, vmname string) (bool, *errs.AppError)
}

type AuthRepositoryDb struct {
	client *mgo.Database
}

func (d AuthRepositoryDb) FindIfNodeExists(resgrp, vmname string) (bool, *errs.AppError) {
	var res domain.Organization
	col := d.client.C("organizations")
	err := col.Find(bson.M{"resourcegroup.resourcegroup_name": resgrp, "resourcegroup.virtual_machine.vm_name": vmname}).One(&res)
	if err == mgo.ErrNotFound {
		return false, nil
	} else {
		return true, nil
	}
}
func (d AuthRepositoryDb) SaveNode(vm domain.VmAll) (*domain.Organization, *errs.AppError) {
	vmLoginCred := domain.VmLogin{
		VmName:     vm.VmName,
		VmUserName: vm.VmUserName,
		VmPassword: vm.VmPassword,
		IpAdd:      vm.IpAdd,
	}
	resGrp := domain.ResourceGroup{
		Name:   vm.ResGrpName,
		Region: vm.Region,
		LoginDet: []domain.VmLogin{
			vmLoginCred,
		},
	}
	org := domain.Organization{
		OrgName:       vm.OrgName,
		ResourceGroup: []domain.ResourceGroup{resGrp},
	}
	loc := domain.Location{
		Region: vm.Region,
	}
	var resLoc domain.Location

	col := d.client.C("organizations")
	col2 := d.client.C("resourcegroup_locations")
	var res domain.Organization

	err0 := col2.Find(bson.M{"region": loc.Region}).One(&resLoc)
	if err0 == mgo.ErrNotFound {
		err0 = col2.Insert(&loc)
		if err0 != nil {
			logger.Error("error while inserting new resourcegroup location" + err0.Error())
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
	}

	err := col.Find(bson.M{"org_name": vm.OrgName}).One(&res)
	if err == mgo.ErrNotFound {
		err = col.Insert(&org)
		if err != nil {
			logger.Error("error while inserting new resource group" + err.Error())
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
		err = col2.Insert(&loc)
		if err != nil {
			logger.Error("error while inserting new region" + err.Error())
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
	} else {
		err = col.Find(bson.M{"org_name": vm.OrgName, "resourcegroup.resourcegroup_name": resGrp.Name}).One(&res)
		if err == mgo.ErrNotFound {
			pushQuery := bson.M{"resourcegroup": resGrp}
			err1 := col.Update(bson.M{"org_name": vm.OrgName}, bson.M{"$addToSet": pushQuery})
			if err1 != nil {
				logger.Error("error while updating resource group array" + err1.Error())
				return nil, errs.NewUnexpectedError("unexpected DB error")
			}

		} else {
			pushQuery := bson.M{"resourcegroup.$.virtual_machine": vmLoginCred}
			err1 := col.Update(bson.M{"org_name": vm.OrgName, "resourcegroup.resourcegroup_name": resGrp.Name}, bson.M{"$addToSet": pushQuery})
			if err1 != nil {
				logger.Error("error while updating resource group" + err1.Error())
				return nil, errs.NewUnexpectedError("unexpected DB error")
			}
		}
	}
	err = col.Find(bson.M{"org_name": vm.OrgName}).One(&res)
	if err != nil {
		logger.Error("error while fetching organization" + err.Error())
		return nil, errs.NewUnexpectedError("unexpected DB error")
	}
	return &res, nil
}
func (d AuthRepositoryDb) ResetPassword(evpw, url_email, email, newpw, confirmpw string) *errs.AppError {
	email_verf_store := d.client.C("verification_code")

	if newpw != confirmpw {
		return errs.NewUnexpectedError("passwords doesn't match")
	}
	if utils.Check(email, url_email) == false {
		return errs.NewUnexpectedError("invalid email")
	}

	var verif_code_struct domain.PwResetParams

	err := email_verf_store.Find(bson.M{"email": email, "code": evpw}).One(&verif_code_struct)
	if err != nil {
		if err == mgo.ErrNotFound {
			logger.Error("Incorrect hash" + err.Error())
			return errs.NewUnexpectedError("incorrect url hash")
		}
		logger.Error("Error while querying email verification code: " + err.Error())
		return errs.NewUnexpectedError("unexpected db error")
	}
	if time.Now().After(verif_code_struct.Timeout) {
		return errs.NewUnexpectedError("timeout has expired")
	} else {
		if evpw == verif_code_struct.Code {
			hash, err := bcrypt.GenerateFromPassword([]byte(confirmpw), 10)
			if err != nil {
				logger.Error("error while encrypting new password :" + err.Error())
				return errs.NewUnexpectedError("unexpected DB error")
			}
			confirmpw = string(hash)

			users := d.client.C("users")
			err1 := users.Update(bson.M{"email": email}, bson.M{"$set": bson.M{"password": confirmpw}})
			if err1 != nil {
				logger.Error("error while updating new password : " + err1.Error())
				return errs.NewUnexpectedError("unexpected database error")
			}
		} else {
			return errs.NewUnexpectedError("email verification hash in url doesn't match ")
		}
	}
	return nil
}
func (d AuthRepositoryDb) RecoverEmail(email string) *errs.AppError {
	now := time.Now()
	timeout := now.Add(time.Minute * 45)
	ifEx, _ := d.findIfEmailExists(email)
	if ifEx == true {
		newCode := utils.GenerateCode()

		strCode := strconv.Itoa(newCode)
		code_hash := utils.EncryptStr(strCode)
		fmt.Println(strCode)

		email_verif_code := d.client.C("verification_code")

		err := email_verif_code.Update(bson.M{"email": email}, bson.M{"$set": bson.M{"code": code_hash, "timeout": timeout}})
		if err != nil {
			logger.Error("error while updating email verification code : " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
		err2 := utils.SendEmail(code_hash, email, "")
		if err2 != nil {
			logger.Error("error while sending email - " + err2.Message)
			return errs.NewUnexpectedError("unexpected DB error")
		}
		return nil
	} else {

		return errs.NewUnexpectedError("email not found")
	}
}
func (d AuthRepositoryDb) findIfEmailExists(email string) (bool, *errs.AppError) {
	users := d.client.C("users")

	err := users.Find(bson.M{"email": email}).One(&users)
	if err == mgo.ErrNotFound {
		return false, nil
	} else {
		return true, errs.NewUnexpectedError("email not found")
	}
}
func (d AuthRepositoryDb) VerifyMobileNo(verify domain.VerifyMobile) *errs.AppError {
	ifSame := d.CheckIfEqualSmsCode(verify.Code, verify.Contact_No)

	if ifSame == true {
		users := d.client.C("users")

		err := users.Update(bson.M{"prof.contact_no": verify.Contact_No}, bson.M{"$set": bson.M{"mobile_verified": true}})

		if err != nil {
			logger.Error("error while updating mobile verification status : " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
		return nil
	} else {
		return errs.NewUnexpectedError("mobile verification code doesn't match")
	}
}
func (d AuthRepositoryDb) VerifyEmail(verify domain.VerifyEmail) *errs.AppError {

	ifSame := d.CheckIfEqualEmailCode(verify.Code, verify.Email)

	if ifSame == true {
		users := d.client.C("users")

		err := users.Update(bson.M{"email": verify.Email}, bson.M{"$set": bson.M{"email_verified": true, "activated": true}})

		if err != nil {
			logger.Error("error while updateing email verification status : " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
		return nil
	} else {
		return errs.NewUnexpectedError("verification code doesn't match")
	}
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) *errs.AppError {
	var rtoken domain.Refresh_Token
	refresh_tokens := d.client.C("refresh_token_store")
	err0 := refresh_tokens.Find(bson.M{"token": refreshToken}).One(&rtoken)
	if err0 != nil {
		if err0 == mgo.ErrNotFound {
			return errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err0.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
	}
	return nil
}

func (d AuthRepositoryDb) SaveUser(user domain.User) (*domain.User, *errs.AppError) {
	users := d.client.C("users")

	err0 := users.Find(bson.M{"email": user.Email}).One(&user)
	if err0 == mgo.ErrNotFound {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			logger.Error("error while encrypting password :" + err.Error())
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
		user.Password = string(hash)

		err = users.Insert(&user)
		if err != nil {
			logger.Error("error while inserting user into DB :" + err.Error())
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
		user.Password = ""

		code := utils.GenerateCode()
		strCode := strconv.Itoa(code)
		code_to_store := strCode
		code_hash, errr := bcrypt.GenerateFromPassword([]byte(strCode), 10)
		if errr != nil {
			logger.Error("error while encrypting code :" + errr.Error())
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
		fmt.Println(code_hash)
		code_hashStr := string(code_hash)
		fmt.Println(code_hashStr)
		err1 := d.SaveEmailVerificationCode(code_hashStr, user.Email)
		if err1 != nil {
			logger.Error("error while inserting code:" + err1.Message)
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
		err2 := utils.SendEmail(code_to_store, user.Email, user.Prof.Firstame)
		if err2 != nil {
			logger.Error("error while sending email - " + err2.Message)
			return nil, errs.NewUnexpectedError("unexpected DB error")
		}
		return &user, nil
	} else {
		return nil, errs.NewUnexpectedError("user has already been registered under the given email")
	}
}

func (d AuthRepositoryDb) Login(email, password string) (*domain.Login, *errs.AppError) {
	var login domain.Login
	user := d.client.C("users")

	err := user.Find(bson.M{"email": email}).One(&login)

	if err != nil {
		logger.Error("error while querying email : ")
		return nil, errs.NewUnexpectedError("invalid email" + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))

	if err != nil {
		logger.Error("error while querying password : " + err.Error())
		return nil, errs.NewUnexpectedError("invalid password")
	}
	if login.IsEmailVerified == false {
		logger.Error("email has not confirmed yet")
		return nil, errs.NewUnexpectedError("email has not confirmed yet")
	}

	code := utils.GenerateCode()
	strCode := strconv.Itoa(code)
	code_hash, errr := bcrypt.GenerateFromPassword([]byte(strCode), 10)
	if errr != nil {
		logger.Error("error while encrypting code :" + errr.Error())
		return nil, errs.NewUnexpectedError("unexpected DB error")
	}

	strCode = string(code_hash)
	err1 := d.SaveMobileVerificationCode(strCode, login.Prof.Contact_No)
	if err1 != nil {
		logger.Error("error while inserting sms verification code:" + err.Error())
		return nil, errs.NewUnexpectedError("unexpected DB error")
	}

	err2 := utils.SendSms(login.Prof.Contact_No, code)
	if err2 != nil {
		return nil, errs.NewUnexpectedError("unexpected DB error")
	}
	return &login, nil
}
func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(authToken domain.AuthToken) (string, *errs.AppError) {
	var appErr *errs.AppError
	var refreshToken string
	var time string = time.Now().Format("2006-01-02 15:04:05")
	if refreshToken, appErr = authToken.NewRefreshToken(); appErr != nil {
		return "", appErr
	}

	rtoken := domain.Refresh_Token{
		Token:     refreshToken,
		TimeStamp: time,
	}
	token := d.client.C("refresh_token_store")

	err := token.Insert(rtoken)
	if err != nil {
		logger.Error("error while inserting refresh token into DB :" + err.Error())
		return "", errs.NewUnexpectedError("unexpected DB error")
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) SaveMobileVerificationCode(code string, mobile string) *errs.AppError {

	code_struct := domain.VerifyMobile{
		Contact_No: mobile,
		Code:       code,
	}

	code_store := d.client.C("mobile_verification_code")
	err := code_store.Insert(code_struct)
	if err != nil {
		logger.Error("error while inserting mobile verification code into DB :" + err.Error())
		return errs.NewUnexpectedError("unexpected DB error")
	}
	return nil
}
func (d AuthRepositoryDb) SaveEmailVerificationCode(code string, em string) *errs.AppError {
	code_struct := domain.VerifyEmail{
		Email: em,
		Code:  code,
	}

	code_store := d.client.C("verification_code")
	err := code_store.Insert(code_struct)
	if err != nil {
		logger.Error("error while inserting email verification code into DB :" + err.Error())
		return errs.NewUnexpectedError("unexpected DB error")
	}

	return nil
}
func (d AuthRepositoryDb) CheckIfEqualEmailCode(code string, email string) bool {
	var req domain.VerifyEmail

	code_store := d.client.C("verification_code")
	err := code_store.Find(bson.M{"email": email}).One(&req)
	if err != nil {
		logger.Error("Error while querying verification code: " + err.Error())
		return false
	}
	err1 := bcrypt.CompareHashAndPassword([]byte(req.Code), []byte(code))
	if err1 != nil {
		logger.Error("error while comparing verification codes : " + err.Error())
		return false
	}

	return true
}
func (d AuthRepositoryDb) CheckIfEqualSmsCode(code string, mobile string) bool {
	var req domain.VerifyMobile

	code_store := d.client.C("mobile_verification_code")
	err := code_store.Find(bson.M{"contact_no": mobile}).One(&req)

	if err != nil {
		logger.Error("Error while querying mobile verification code: " + err.Error())
		return false
	}
	err1 := bcrypt.CompareHashAndPassword([]byte(req.Code), []byte(code))
	if err1 != nil {
		logger.Error("error while comparing verification codes : " + err1.Error())
		return false
	}
	return true
}
func NewAuthRepository(client *mgo.Database) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
