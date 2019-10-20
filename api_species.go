package firstAssignment

import (
	"encoding/json"
	"net/http"
	"strings"
	"bytes"
	"strconv"
)



type YearHolder struct{													// Struct for species/key/name api
	Year									string
}



var structure = new(ResultList)									// Structure for JSON result field that holds array of species
var DBs = SpeciesDB{}														// WHere all species are stored



func fetchSpeciesJSON(url string) ([]Species, error) {	// Fetches specific json reauested
	resp, err := http.Get(url)										// Gets information from url
	if err != nil {																// Checks for error
		return nil, err															// Returns error i found
	}
	defer resp.Body.Close()												// Opens html body
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &structure); err != nil {	// Copies the json content into structure
		return nil, err															// Returns error if found
	}
	for idx, row := range structure.AllSpecies {	// For each species
		if idx == 0 {																// Skip first slot
			continue
		}
		url := "http://api.gbif.org/v1/species/"		// Open species/key/name for each reteieved species to find date
		url = url + strconv.Itoa(int(row.Key)) + "/name"	// Gets species key and uses for url
		resp, err := http.Get(url)									// Fetches content from url
		if err != nil {
			return nil, err														// Error handling
		}
		var yearFetch YearHolder										// Temp for result containing date
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respByte := buf.Bytes()
		json.Unmarshal(respByte, &yearFetch)				// Copies json into variable
		row.Year = yearFetch.Year										// Assigns the fetched year to species
		structure.AllSpecies[idx] = row							// Places species in map
	}

	return structure.AllSpecies, nil							// Returns map with structures
}





//-------------------------------------------------------------

																								// Prints all species
func replyWithAlls(w http.ResponseWriter, DB speciesStorage) {
																								// Url showing a few Species just for if user doesnt write key in url
		url := "http://api.gbif.org/v1/species?offset=20000&limit=25"
		speciesList, err := fetchSpeciesJSON(url)		// Fetches using function above
		if err != nil {															// Error handling
			http.Error(w, "Service could not be accessed", http.StatusServiceUnavailable)
		}

		for idx, row := range speciesList {					// For each species
		  if idx == 0 {
		  	continue
	  	}
      DBs.Add(row)															// Put in main Species struct
    }


		if DB.Count() == 0 {												// If none are in the map, write out empty JSON
			json.NewEncoder(w).Encode([]Species{})
		} else {																		// If there are species stored
			a := make([]Species, 0, DB.Count())				// Fetch them all
			for _, s := range DB.GetAll() {
				a = append(a, s)
			}
			json.NewEncoder(w).Encode(a)							// And print as JSON to website
	}
}





func replyWithSpecifics(w http.ResponseWriter, DB speciesStorage, id string) {
		var singleFetch Species														// Temp variable for fetching species
		url := "http://api.gbif.org/v1/species/" + id			// Fetches species with ID in url
		resp, err := http.Get(url)												// Fetches url
		if err != nil {
			http.Error(w, "Species Service Unavailable", http.StatusServiceUnavailable)
			DN.TestApi("species")														// Returns error code 503 if not working
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respByte := buf.Bytes()
		json.Unmarshal(respByte, &singleFetch)						// Fetches specific species

		DBs.Add(singleFetch)															// Stores it
		idINT, err := strconv.ParseUint(id, 10, 64)				// Casts it to unit64
		if err != nil {
			http.Error(w, "Species Service Unavailable", http.StatusServiceUnavailable)
			DN.TestApi("species")
		}
		s, ok := DBs.Get(idINT)														// FEtches species with key = id
		if !ok {
			http.Error(w, "Species not found", http.StatusNotFound)
			return																					// status not found if species is not in the map
		}
		json.NewEncoder(w).Encode(s)											// Prints relevant info as JSON to website
		}


func HandlerSpecies(w http.ResponseWriter, r *http.Request) {
		http.Header.Add(w.Header(), "content-type", "application/json")	// Sets web application type
		parts := strings.Split(r.URL.Path, "/")						// Splits url into parts divided by /
		if len(parts) == 6 || parts[1] != "conservation"{ // If not too long or doesnt contain required key
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if parts[4] == "" {																// If 4th part is empty reply with a standard few species
			replyWithAlls(w, &DBs)
		} else {																					// Otherwise, print a specific species
			replyWithSpecifics(w, &DBs, parts[4])					
		}
}
