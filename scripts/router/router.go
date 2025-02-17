package router

import (
	"net/http"

	"github.com/golang-mysql/scripts/image"
	"github.com/golang-mysql/scripts/image_v2"
	"github.com/golang-mysql/scripts/durasirawat"
	"github.com/golang-mysql/scripts/utils"
	"github.com/gorilla/mux"
)

var Index = newIndexHandler()

func newIndexHandler() http.Handler {
	router := mux.NewRouter()
	router.Use(commonMiddleware)

	// Api Resize Image
	router.HandleFunc("/api/v1/image", image.Upload).Methods("POST")
	// Resize image_v2
	router.HandleFunc("/api/v2/image", image_v2.Upload).Methods("POST")

	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("/app/images/"))))

	router.HandleFunc("/healthz", utils.Healthz).Methods("GET")

	router.HandleFunc("/durasi", durasirawat.DurasiHandler).Methods("GET")
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
