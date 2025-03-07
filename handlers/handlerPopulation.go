package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Conv -Struct to convert iso2 to iso3 country code
type Conv struct {
	Iso3 string `json:"cca3"`
}

// PopulationTemp - hold the temporary structure for the fetched population data
type PopulationTemp struct {
	Data struct {
		PopulationCounts []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

// FinalPopulation - Holds the filtered/formatted population values and their avg(mean)
type FinalPopulation struct {
	Mean   int `json:"mean"`
	Values []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	} `json:"values"`
}

// HandlerPopulation handles incoming http requests and retrieves the population data
func HandlerPopulation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		//iso - for the two-letter-country-code
		iso := r.PathValue("two_letter_country_code")
		if len(iso) != 2 { //Restricting to only 2 letters (iso2 format)
			http.Error(w, "Invalid ISO, expecting a 2 letter country code =)"+LINEBREAK+"TRY AGAIN", http.StatusBadRequest)
			return
		}

		//Gets iso3 countrycode from the iso2 above, returns if there is an error.
		iso3, err := ConvertISO(w, iso)
		if err != nil {
			return
		}

		var start, end int
		// Get the limit query parameter, used to filter population data by year range
		limit := r.URL.Query().Get("limit")
		if limit == "" {
			start = 0               //Default: this is where Year starts from
			end = time.Now().Year() //End at the current year
		} else {
			years := strings.Split(limit, "-") //Expecting format: XXXX-YYYY
			if len(years) != 2 {               //We limit the array (that gets created in string.split) to 2
				http.Error(w, "Invalid input for limit"+LINEBREAK+"Use `?limit=XXXX-YYYY`", http.StatusBadRequest)
				return
			}
			// Parse the start year from string to int
			s, startErr := strconv.Atoi(years[0])
			if startErr != nil {
				http.Error(w, "Start (XXXX) must be a number."+LINEBREAK+"TRY AGAIN =)", http.StatusBadRequest)
				return
			} else {
				start = s
			}
			// Parse the end year from string to int
			e, endErr := strconv.Atoi(years[1])
			if endErr != nil {
				http.Error(w, "End (YYYY) must be a number"+LINEBREAK+"TRY AGAIN =)", http.StatusBadRequest)
				return
			} else {
				end = e
			}
		}

		//Creating a temp variable for the finalPopulation struct
		var finalPopTemp FinalPopulation

		//Fetch population data and store it in the temp struct above.
		err2 := GETPopulation(w, iso3, start, end, &finalPopTemp)
		if err2 != nil {
			return
		}

		//Encode the final population data as JSON and sent it as a response
		errJsn := json.NewEncoder(w).Encode(finalPopTemp)
		if errJsn != nil {
			http.Error(w, errJsn.Error(), http.StatusInternalServerError)
			return
		}
	} else { //Telling our user that they can only use GET, nothing else!
		http.Error(w, "Method not allowed. ", http.StatusMethodNotAllowed)
	}
}

// GETPopulation fetched population data for a country within a given year range
func GETPopulation(w http.ResponseWriter, iso3 string, start int, end int, finalPopTemp *FinalPopulation) error {
	w.Header().Set("Content-Type", `application/json`)

	//Variable URL for the API request
	var popURL = COUNTRYURL + "countries/population"

	//Preparing a payload with the iso3 country code
	payloadData := map[string]string{"iso3": iso3} //Preparing a map with a single key set as iso3
	payloadBytes, err := json.Marshal(payloadData) //Converts payloadData map into a JSON slice
	if err != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}
	//makes payload readable to use as the request body
	payload := strings.NewReader(string(payloadBytes))

	//Initialize the http client
	client := &http.Client{}
	defer client.CloseIdleConnections() //defer to make sure it closes when function ends

	// Create a POST request to fetch population data
	req, err := http.NewRequest("POST", popURL, payload)
	if err != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}
	req.Header.Set("Content-Type", "application/json")

	//Execute the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	defer func(Body io.ReadCloser) { //another defer to make sure the body closes when function ends
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check if API returned an error
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API error %s", resp.Status), resp.StatusCode)
		return errors.New("")
	}

	// Read the response
	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	//Parse the JSON response into the temporary population struct
	var populationData PopulationTemp
	errJSONs := json.Unmarshal(body, &populationData)
	if errJSONs != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return errors.New("")
	}

	//Iterate through population data and filter by the given year range.
	var meanSum int
	for i := 0; i < len(populationData.Data.PopulationCounts); i++ {
		//iterate and appends if the input is >= the XXXX and <= the YYYY (XXXX-YYYY)
		if populationData.Data.PopulationCounts[i].Year >= start && populationData.Data.PopulationCounts[i].Year <= end {
			finalPopTemp.Values = append(finalPopTemp.Values, populationData.Data.PopulationCounts[i])

			//meanSum contains all the population added together within the year range that is requested.
			meanSum += populationData.Data.PopulationCounts[i].Value
		}
	}
	// If the struct is 0, we set the mean (avg) to 0
	if len(finalPopTemp.Values) == 0 {
		finalPopTemp.Mean = 0
	} else { //Otherwise we calculate the mean(avg) value of the population by dividing by the amount of years we retrieve
		finalPopTemp.Mean = meanSum / len(finalPopTemp.Values)
	}
	return nil
}

// ConvertISO converts a 2letter country code (iso2) into a 3letter country code(iso3)
func ConvertISO(w http.ResponseWriter, iso string) (string, error) {

	w.Header().Set("Content-Type", `application/json`)

	// Create an HTTP client to make requests
	client := &http.Client{}
	defer client.CloseIdleConnections() // Cleanup: Close any idle connections when we're done

	//Constructing the API request for the RESTURL and getting the cca3 (iso3)
	var isoURL = RESTURL + "alpha/" + iso + "?fields=cca3"

	//Creating a GET request to fetch the iso3 country code
	req, errREQ := http.NewRequest("GET", isoURL, nil)
	if errREQ != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return "", errors.New("")
	}

	//Execute the request!
	res, errDO := client.Do(req)
	if errDO != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return "", errors.New("")
	}
	//Handle case where the given iso2 code doesn't exist
	if res.StatusCode == http.StatusNotFound {
		http.Error(w, "Sorry dude, or dudette. But your iso code is very very bad! Because it doesnt exist!"+LINEBREAK+"TRY AGAIN", http.StatusNotFound)
		return "", errors.New("")
	}

	//Read the API response
	body, errRead := io.ReadAll(res.Body)
	if errRead != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return "", errors.New("")
	}

	defer func(Body io.ReadCloser) { //defer to make sure it closes when function ends
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	// Parse the JSON response into the Converter struct
	var convData Conv
	errJSON := json.Unmarshal(body, &convData)
	if errJSON != nil {
		http.Error(w, "Internal Server Error :(", http.StatusInternalServerError)
		return "", errors.New("")
	}
	//returns the converted iso3 country code
	return convData.Iso3, nil

}
