package main

import (
	"log"
	"net/http"
)

func main() {
	ConnectDatabase()

	router := http.NewServeMux()

	// Creating a router

	// Register Organisation endpoint
	router.HandleFunc("/organisations", OrganisationHandler)

	// Register User endpoint
	router.HandleFunc("/users", UserHandler)

	log.Println("Server Started on Port 8080")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}

}
