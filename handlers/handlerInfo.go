package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// GETCountryInfo struct to get and display the retrieved data on our API
type GETCountryInfo struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    []string          `json:"capital"`
	Cities     []string          `json:"cities"`
}

// JustCities struct to retrieve city data from COUNTRY-API
type JustCities struct {
	Cities []string `json:"data"`
}

// HandlerInfo - this is the main handler that deals with country information
func HandlerInfo(w http.ResponseWriter, r *http.Request) {

	// Country - Global variable for the GET-struct, where all the fetched data is stored
	country := GETCountryInfo{}

	//iso - for the two-letter-country-code
	iso := r.PathValue("two_letter_country_code")
	if len(iso) != 2 { //Restricting to only 2 letters
		http.Error(w, "Invalid iso. Have to be a two letter iso (country code)", http.StatusBadRequest)
		return
	}
	//Fetch country details and store them in country struct
	err := GetCountry(w, r, iso, &country)
	if err != nil {
		return
	}

	// Get the limit query parameter, used to limit city results
	limit := r.URL.Query().Get("limit")

	// If the limit for cities is empty, default set to 10
	if limit == "" {
		err2 := FetchCities(w, r, &country, iso, 10)
		if err2 != nil {
			return
		}
	} else {
		//Convert limit to an int, and making sure its a number
		limitInt, limitErr := strconv.Atoi(limit)
		if limitErr != nil {
			http.Error(w, `Hello dude, or dudette. Limit has to be a numeric value`, http.StatusInternalServerError)
			return
		}
		if limitInt <= 0 { //Making sure that negative numbers aren't allowed
			http.Error(w, "Limit must be a positive number =)", http.StatusBadRequest)
			return
		}
		//Calling FetchCities() to get city data for country and their limit
		err3 := FetchCities(w, r, &country, iso, limitInt)
		if err3 != nil {
			return
		}
	}

	//Encode the final country struct as JSON, sending the response to my API =)
	errJsn := json.NewEncoder(w).Encode(country)
	if errJsn != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return
	}
}

// GetCountry - retrieves data from the REST API with GET-method
func GetCountry(w http.ResponseWriter, r *http.Request, iso string, c *GETCountryInfo) error {
	if r.Method == http.MethodGet { //Making sure the user can only use GET

		//Setting response content type
		w.Header().Set("Content-Type", `application/json`)

		// Create an HTTP client to make request
		client := &http.Client{}
		defer client.CloseIdleConnections() // Cleanup: Close any idle connections when we're done

		//Variable for the URL we want to use (using the iso)
		var everythingURL = RESTURL + "alpha/" + iso

		////Making a new http request using GET (Using my URL variable ^)
		req, errREQ := http.NewRequest("GET", everythingURL, nil)
		if errREQ != nil {
			http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
			return errors.New("")
		}

		//Sending the request and receive the response
		res, errDO := client.Do(req)
		if errDO != nil {
			http.Error(w, errDO.Error(), http.StatusInternalServerError)
			return errors.New("")
		}
		//If the API responds with 404, the iso input probably doesn't exist
		if res.StatusCode == http.StatusNotFound {
			http.Error(w, "Sorry dude, or dudette. But your iso code is very very bad! Because it doesnt exist!"+LINEBREAK+"TRY AGAIN", http.StatusNotFound)
			return errors.New("")
		}
		//Reading the response (body)
		body, errRead := io.ReadAll(res.Body)
		if errRead != nil {
			http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
			return errors.New("")
		}

		//Closing the response body to avoid memory leaks.
		//..because it might not close at the end of the function
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
			}
		}(res.Body)

		// Expecting an array and unmarshal into a slice
		var countryData []GETCountryInfo
		errJSON := json.Unmarshal(body, &countryData)
		if errJSON != nil {
			http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
			return errors.New("")
		}

		// Ensure that if we countryData is empty, we get a sad response and error.
		if len(countryData) == 0 {
			http.Error(w, "No country data found", http.StatusNotFound)
			return errors.New("")
		}

		// Assign the first country from the array
		*c = countryData[0]

	} else { //Telling our user that they can only use GET, nothing else!
		http.Error(w, "Method not allowed. ONLY USE GET!", http.StatusMethodNotAllowed)
	}
	return nil
}

// FetchCities - getting cities information from the second API (COUNTRYURL)
func FetchCities(w http.ResponseWriter, _ *http.Request, c *GETCountryInfo, iso string, limit int) error {

	//Debug info for me
	fmt.Println("\nFetching cities for:", iso, "with limit:", limit)

	//Variable for the URL we want to use (for fetching city data)
	var citiesURL = COUNTRYURL + "countries/cities"

	//Creating a payload with the country's iso2 code (NO, SE, DK)
	payloadData := map[string]string{"iso2": iso}
	payloadBytes, err := json.Marshal(payloadData) //converting to JSON
	if err != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	//Converting the payload to a readable format :D
	payload := strings.NewReader(string(payloadBytes))

	//Create an HTTP client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Create a POST request to fetch cities
	req, err := http.NewRequest("POST", citiesURL, payload)
	if err != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	//setting request headers
	req.Header.Set("Content-Type", "application/json")

	//Execute the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	defer func(Body io.ReadCloser) { //defer to make sure it closes when function ends
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	// Check if API returned an error
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API error: %s", resp.Status), resp.StatusCode)
		return errors.New("")
	}

	// Read the response
	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	//Parse the JSON response
	var temp JustCities
	errJSONs := json.Unmarshal(body, &temp)
	if errJSONs != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}
	//Sorting cities in alphabetically
	sort.Strings(temp.Cities)

	//Slicing cities into the limit parameter
	c.Cities = temp.Cities[:limit]

	return nil
}
