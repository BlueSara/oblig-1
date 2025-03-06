package temp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	. "log"    // Import log functions without prefix
	"net/http" // Used for handling HTTP requests
	"strconv"
)

const COUNTRYURL = "http://129.241.150.113:3500/api/v0.1/"
const RESTURL = "http://129.241.150.113:8080/v3.1/"

// Struct to hold API response data about a country
type GETCountryInfo struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"png"`
	Capital    []string          `json:"capital"`
	Cities     []string          `json:"data"`
}

// Struct for formatted country information to return to the client
type CountryInfo struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"language"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flags"`
	Capital    string            `json:"capital"`
	Cities     []string          `json:"cities"`
}

// HTTP handler function for fetching country info
func HandlerInfo(w http.ResponseWriter, r *http.Request) {

	// Only allow GET requests
	switch r.Method {
	case http.MethodGet:

		// Extract country code from the request path
		countryCode := r.PathValue("two_letter_country_code")

		// Get the 'limit' query parameter, used to limit city results
		limit := r.URL.Query().Get("limit")

		// Convert limit to an integer if it's provided
		limitInt, err := strconv.Atoi(limit)
		if err != nil && limit != "" {
			Printf("failed to convert limit parameter to int: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Build the API request URL for country information
		reqURL := RESTURL + fmt.Sprintf("alpha/%s?fields=name,continents,population,languages,borders,flags,capital", countryCode)

		// Make a GET request to the API
		resp, err := http.Get(reqURL)
		if err != nil {
			Printf("failed to fetch data from API: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Ensure the response body is closed after function exits
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				Printf("failed to close response body: %v", err)
			}
		}(resp.Body)

		// Read the response body (but from the wrong source `r.Body` instead of `resp.Body`)
		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			Printf("failed to read response body: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Parse JSON response into GETCountryInfo struct
		var country GETCountryInfo
		if err := json.Unmarshal(respBody, &country); err != nil {
			Printf("failed to unmarshal response body: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Fetch cities for the country
		cities, err := getCities(country.Name.Common, limitInt)
		if err != nil {
			Printf("failed to get cities: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Format the response to match the CountryInfo struct
		countryInfo := CountryInfo{
			Name:       country.Name.Common,
			Continents: country.Continents,
			Population: country.Population,
			Languages:  country.Languages,
			Borders:    country.Borders,
			Flag:       country.Flag,
			Capital:    country.Capital[0],
			Cities:     cities,
		}

		// Set response content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Send the response as JSON
		if err := json.NewEncoder(w).Encode(countryInfo); err != nil {
			Printf("failed to encode response body: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// Function to get city names for a given country
func getCities(country string, limit int) ([]string, error) {
	// Prepare the request body with the country name
	requestBody, _ := json.Marshal(map[string]string{
		"country": country,
	})

	// Create a POST request to fetch cities
	req, err := http.NewRequest("POST", COUNTRYURL+"countries/cities", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set content type as JSON
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Ensure response body is closed after function exits
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Println("Error closing body:", err)
		}
	}(resp.Body)

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		Println("Error reading response body:", err)
	}

	// Parse the response into the CountryInfo struct
	var citiesResp CountryInfo
	if err := json.Unmarshal(respBody, &citiesResp); err != nil {
		Println("Error parsing JSON:", err)
	}

	// If a limit is set, truncate the list of cities
	if limit != 0 && len(citiesResp.Cities) > limit {
		citiesResp.Cities = citiesResp.Cities[:limit]
	}

	return citiesResp.Cities, nil
}
