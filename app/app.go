package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/niluwats/bethel_dashboard/dbhandler"
	"github.com/niluwats/bethel_dashboard/service"
	"gopkg.in/mgo.v2"
)

func Start() {
	router := mux.NewRouter()
	dbClient := getDbClient()
	authRepo := dbhandler.NewAuthRepository(dbClient)
	authHand := AuthHandlers{service.NewAuthService(authRepo)}

	router.HandleFunc("/auth/users", authHand.newUser).Methods(http.MethodPost).Name("NewUser")
	router.HandleFunc("/auth/users/verifyemail", authHand.verifyEmail).Methods(http.MethodPost)
	router.HandleFunc("/auth/users/verifymobile", authHand.verifyMobile).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authHand.login).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", authHand.refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/recoveraccount", authHand.recoverAccount).Methods(http.MethodPost)
	router.HandleFunc("/auth/resetpassword/{email}/{evpw}", authHand.resetPassword).Methods(http.MethodPost)

	address := "localhost"
	port := "8000"
	server := fmt.Sprintf("%s:%s", address, port)
	log.Fatal(http.ListenAndServe(server, router))
}
func getDbClient() *mgo.Database {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}
	db := session.DB("bethel_dashboard")
	return db
}
