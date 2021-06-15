package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/niluwats/bethel_dashboard/dbhandler"
	"github.com/niluwats/bethel_dashboard/service"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func Start() {
	router := mux.NewRouter()
	dbClient := getDbClient()
	authRepo := dbhandler.NewAuthRepository(dbClient)
	authHand := AuthHandlers{service.NewAuthService(authRepo)}

	router.HandleFunc("/v1/auth/users", authHand.newUser).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/users/verifyemail/{email}/{evpw}", authHand.verifyEmail).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/users/verifymobile", authHand.verifyMobile).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/login", authHand.login).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/users/sendsms", authHand.getMbVerificationCode).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/refresh", authHand.refresh).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/recoveraccount", authHand.recoverAccount).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/resetpassword/{email}/{evpw}", authHand.resetPassword).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/nodes", authHand.createNode).Methods(http.MethodPost).Name("NewNode")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
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
