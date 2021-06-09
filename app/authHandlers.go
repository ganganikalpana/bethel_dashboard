package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/niluwats/bethel_dashboard/dto"
	"github.com/niluwats/bethel_dashboard/logger"
	"github.com/niluwats/bethel_dashboard/service"
)

type AuthHandlers struct {
	service service.AuthService
}

func (h AuthHandlers) createNode(w http.ResponseWriter, r *http.Request){
	var req dto.NewNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Error while decoding create node request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		nodes,appErr := h.service.CreateNode(req)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, nodes)
		}
	}

}

func (h AuthHandlers) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.PasswordReset
	vars := mux.Vars(r)
	evpw := vars["evpw"]
	url_email := vars["email"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Error while decoding password reset request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		appErr := h.service.ResetPassword(req, evpw, url_email)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, "Reset password successfully")
		}
	}
}
func (h AuthHandlers) recoverAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.RecoverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Error while decoding recover request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		appErr := h.service.Recover(req)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, "Please check your inbox")
		}
	}
}
func (h AuthHandlers) refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		logger.Error("Error while decoding refresh token request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Refresh(refreshRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}
func (h AuthHandlers) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var request dto.EmailVerifyCode
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		appErr := h.service.VerifyEmail(request)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.Message)
		} else {
			writeResponse(w, http.StatusCreated, "verified email successfully")
		}
	}
}
func (h AuthHandlers) verifyMobile(w http.ResponseWriter, r *http.Request) {
	var request dto.MobileVerifyCode
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		appErr := h.service.VerifyMobile(request)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.Message)
		} else {
			writeResponse(w, http.StatusCreated, "verified mobile number successfully")
		}
	}
}
func (h AuthHandlers) newUser(w http.ResponseWriter, r *http.Request) {
	var request dto.NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		user, appErr := h.service.Register(request)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.Message)
		} else {
			writeResponse(w, http.StatusCreated, user)
		}
	}
}
func (h AuthHandlers) login(w http.ResponseWriter, r *http.Request) {
	var request dto.NewLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Login(request)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
