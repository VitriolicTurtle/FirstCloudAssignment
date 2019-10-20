package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"firstAssignment"
)

																						// Default answer for no text after /
func handlerNil(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}



func main() {


	//------------------------------------
																						// Initialises the 4 struct maps:
	firstAssignment.DBc.Init()
	firstAssignment.DBs.Init()
	firstAssignment.DN.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

																						// Runs functions based on url typed
	http.HandleFunc("/", handlerNil)
	http.HandleFunc("/conservation/v1/country/", firstAssignment.HandlerCountry)
	http.HandleFunc("/conservation/v1/species/", firstAssignment.HandlerSpecies)
	http.HandleFunc("/conservation/v1/diag/", firstAssignment.HandlerDiag)
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
