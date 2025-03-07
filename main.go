package main

import (
	"Assignment1/handlers" // We import our own handlers, cuz we need them to do stuff
	"log"
	"net/http" // This thing lets us do internet magic
	// We use this to grab env variables, if any
)

func main() {

	// Make a thing that know where to send requests (router)
	router := http.NewServeMux()

	// When user type "/countryinfo/v1/info/" in browser, we call HandlerInfo etc.
	router.HandleFunc("/countryinfo/v1/info/{two_letter_country_code}", handlers.HandlerInfo)
	router.HandleFunc("/countryinfo/v1/population/{two_letter_country_code}", handlers.HandlerPopulation)
	router.HandleFunc("/countryinfo/v1/status/", handlers.HandlerStatus)

	// Try to start server. If fail, we cry and print error
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Oh no! Server go boom:", err.Error())
	}

}
