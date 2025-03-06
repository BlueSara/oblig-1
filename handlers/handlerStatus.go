package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// A timer to track how long the service has been running
var tiktok = time.Now()

// Struct to hold the status of different APIs, version, and uptime.
type statusAPI struct {
	CountriesNowApi  string `json:"countriesnowapi"`
	RestCountriesApi string `json:"restcountriesapi"`
	Version          string `json:"version"`
	Uptime           int64  `json:"uptime"`
}

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Create an HTTP client to make requests
		client := &http.Client{}
		defer client.CloseIdleConnections() // Cleanup: Close any idle connections when we're done

		// Function to check API status
		checkStatus := func(url string) string {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Println("Error creating request:", err.Error())
				return "NOT OK"
			}

			req.Header.Add("Content-Type", "application/json")

			res, err := client.Do(req)
			if err != nil {
				fmt.Println("Error in response:", err.Error())
				return "NOT OK"
			}
			// Ensures to close the Body after the function returns
			defer func(Body io.ReadCloser) {
				err := Body.Close() //Attempt to close the response body.
				if err != nil {     //Is there an error?
					fmt.Println("Error closing body:", err.Error()) //well, then we gotta tell us about it =)
				}
			}(res.Body) //ensure red.Body is properly closed later when the function exits.

			// If we get a response, return the status code
			return fmt.Sprintf("200 OK")
		}

		// Check both APIs
		countriesNowStatus := checkStatus(COUNTRYURL)
		restCountriesStatus := checkStatus(RESTURL)

		// Prepare the response structure
		response := statusAPI{
			CountriesNowApi:  countriesNowStatus,
			RestCountriesApi: restCountriesStatus,
			Version:          "v1",
			Uptime:           int64(time.Since(tiktok).Seconds()), // Calculate uptime in seconds
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Use json.Encoder to directly write the response struct as JSON
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			// If encoding fails, log the error and return an internal server error response "500"
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			fmt.Println("Error encoding response:", err.Error())
		}
	} else {
		http.Error(w, "Method not allowed. ", http.StatusMethodNotAllowed)
	}
}
