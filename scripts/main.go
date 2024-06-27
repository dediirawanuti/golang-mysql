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
		AllowedMethods: []string{"POST, GET, OPTIONS, PUT, DELETE"}, // Allowing only get, just an example
	})

	srv := &http.Server{
		Addr:    ":8910",
		Handler: c.Handler(router.Index),
	}

	fmt.Println("Port" + srv.Addr)

	srv.ListenAndServe()	

}
