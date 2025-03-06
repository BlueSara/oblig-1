package main

import (
	"Assignment1/handlers" // We import our own handlers, cuz we need them to do stuff
	"fmt"                  // For print error when things go boom
	"log"
	"net/http" // This thing lets us do internet magic
	"os"       // We use this to grab env variables, if any
)

func main() {

	// Port for server. Default is 8080, cuz why not
	port := ":8080"

	// if someone (like cloud host) say "use this port," we listen
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT") // Overwrite default port if needed
	}

	// Make a thing that know where to send requests (router)
	router := http.NewServeMux()

	// When user type "/countryinfo/v1/info/" in browser, we call HandlerInfo etc.
	router.HandleFunc("/countryinfo/v1/info/{two_letter_country_code}", handlers.HandlerInfo)
	router.HandleFunc("/countryinfo/v1/population/{two_letter_country_code}", handlers.HandlerPopulation)
	router.HandleFunc("/countryinfo/v1/status/", handlers.HandlerStatus)

	// Print to console so we know server running (very important)
	fmt.Println("Server is running on http://localhost:" + port)

	// Try to start server. If fail, we cry and print error
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Oh no! Server go boom:", err.Error())
	}

}
