package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/niluwats/bethel_dashboard/dbhandler"
	"github.com/niluwats/bethel_dashboard/service"
	"gopkg.in/mgo.v2"
)

func Start() {
	router := mux.NewRouter()
	router.Use(accessControlMiddleware)
	dbClient := getDbClient()
	authRepo := dbhandler.NewAuthRepository(dbClient)
	authHand := AuthHandlers{service.NewAuthService(authRepo)}

	router.HandleFunc("/auth/users", authHand.newUser).Methods(http.MethodPost)
	router.HandleFunc("/auth/users/verifyemail", authHand.verifyEmail).Methods(http.MethodPost)
	router.HandleFunc("/auth/users/verifymobile", authHand.verifyMobile).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authHand.login).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", authHand.refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/recoveraccount", authHand.recoverAccount).Methods(http.MethodPost)
	router.HandleFunc("/auth/resetpassword/{email}/{evpw}", authHand.resetPassword).Methods(http.MethodPost)
	router.HandleFunc("/auth/nodes", authHand.createNode).Methods(http.MethodPost).Name("NewNode")

	// address := "localhost"
	// port := "8000"
	// server := fmt.Sprintf("%s:%s", address, port)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
func getDbClient() *mgo.Database {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}
	db := session.DB("bethel_dashboard")
	return db
}
