package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/golang-mysql/scripts/router"
)

func main() {

	fmt.Println("run main")
	fmt.Println(os.Getenv("ENVIRONMENT"))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                               // All origins
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}, // Allowing only get, just an example
	})

	srv := &http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: c.Handler(router.Index),
	}

	srv.ListenAndServe()	

}
