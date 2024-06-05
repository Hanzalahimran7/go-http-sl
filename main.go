package main

import (
	"log"
	"net/http"

	"github.com/hanzalahimran7/go-http-sl/app"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	a := app.Initialise()
	a.LoadRoutes()
	if err := http.ListenAndServe(":3000", a.Router); err != nil {
		log.Fatal("Failed to initialise the server")
	}
}
