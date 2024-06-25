package router

import (
	"net/http"

	"github.com/golang-mysql/scripts/utils"
	"github.com/gorilla/mux"
)

var Index = newIndexHandler()

func newIndexHandler() http.Handler {
	router := mux.NewRouter()
	router.Use(commonMiddleware)

	router.HandleFunc("/healthz", utils.Healthz).Methods("GET")
	utils.Cron()

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}