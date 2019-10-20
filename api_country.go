package firstAssignment

import (
	"encoding/json"
	"net/http"
	"strings"
	"bytes"
)


type OccurenceList struct {									// Struct for fetching result field
	Oresult[] 					Occurence `json:"results"`
}

type Occurence struct {											// COntains relevant occurence information
	CountryCode 					string `json:"countryCode"`
	GenericName						string `json:"genericName"`
	SpeciesKey						uint64 `json:"speciesKey"`
}


var OCCstructure = new(OccurenceList)				// Structure for Occurence
var DBc = CountriesDB{}											// Stores countries

																						// Fetches crountries from url
func fetchCountryJSON(url string) ([]Country, error) {
	resp, err := http.Get(url)								// GETs url
	if err != nil {														// If it doesnt work, return error
		return nil, err
	}

	defer resp.Body.Close()
	var countryList []Country
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()										// Copy COuntry information
	if err := json.Unmarshal(respByte, &countryList); err != nil {
		return nil, err													// Return error if it fails
	}

	return countryList, nil										// return array of countries
}

																						// Fetches all countries
func HandlerGetCountry(w http.ResponseWriter) {
	url := "https://restcountries.eu/rest/v2/all"		// SPecify url
	countryList, err := fetchCountryJSON(url)				// Uses fetch function
	if err != nil {														// Error handling if fetching failed
		http.Error(w, "Country Service Unavailable", http.StatusServiceUnavailable)
		DN.TestApi("country")										// Assign error to diagnostics
	}

	for idx, row := range countryList {				// For each country in country list
		if idx == 0 {														// SKip header
			continue
		}

		DBc.Add(row)														// Add in DBc map
	}
}



//-------------------------------------------------------------




func replyWithAllc(w http.ResponseWriter, DB countryStorage) {

		if DB.Count() == 0 {										// If it contains nothing print nothing
			json.NewEncoder(w).Encode([]Country{})
		} else {																// Otherwise
			a := make([]Country, 0, DB.Count())		// make map variable for printing
			for _, s := range DB.GetAll() {				// For each country in DB
				a = append(a, s)										// Copy them to a
			}
			json.NewEncoder(w).Encode(a)					// Display as JSON on browser
		}
	}






																						// Reply with specific country with specified amount of occurrences
func replyWithSpecificc(w http.ResponseWriter, DB countryStorage, id string, limit string) {
	url := "http://api.gbif.org/v1/occurrence/search?country=" + id + "&limit=" + limit
	resp, err := http.Get(url)								// get url above
	if err != nil {														// Error handling for if fetching fails
		http.Error(w, "Occurrence Service Unavailable", http.StatusServiceUnavailable)
		DN.TestApi("occurrence")								// Assign error to diagnostics
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	json.Unmarshal(respByte, &OCCstructure)		// Copy JSON Into map

	for idx, x := range OCCstructure.Oresult{ // For each occurrence
		if idx == 0 {														// Skip header
			continue
		}
		DBc.AssignSpecies(x)										// Assign the species (occurrence) to the specific country using function
	}

	s, ok := DB.Get(id)												// Get country from DBc
	if !ok {																	// CHeck if it exists
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
		json.NewEncoder(w).Encode(s)							// Print as JSON to creen
}





																							// Handles the request
func HandlerCountry(w http.ResponseWriter, r *http.Request) {
		HandlerGetCountry(w)											// Fetches all countries upon initialisation
		http.Header.Add(w.Header(), "content-type", "application/json")	// Assigns content type
		parts := strings.Split(r.URL.Path, "/")		// Splits url into string variables
		var limit string = r.URL.Query().Get("limit")	// Fetches limit query
		if limit == ""{														// If specific limit not requsted:
			limit = "10"														// Set to 10 by default
		}

		if len(parts) == 6 || parts[1] != "conservation" {	// Errorhandling checking essential parts of link
			http.Error(w, "Bad request:", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {												// If 4th section does not contain counry code:
			replyWithAllc(w, &DBc)									// Reply with all
		} else {																	// Otherwise
			replyWithSpecificc(w, &DBc, parts[4], limit)	// Reply with a specific country
		}
}
