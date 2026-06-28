package main

import (
	"log"
	"net/http"
)

func main() {
	ConnectDatabase()

	router := http.NewServerMux()

	router.HandleFunc("/organisations", OrganisationHandler)

	router.HandleFunc("/users", UserHandler)

	log.Println("Server Started on Port 8080")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}

}
