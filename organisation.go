package main

import (
	"encoding/json"
	"net/http"
)

type Organisation struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

var organisations []Organisation

func OrganisationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		var org Organisation

		err := json.NewDecoder(r.Body).Decode(&org)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if org.Name == "" {
			http.Error(w, "Organization Name Required", http.StatusBadRequest)
			return
		}

		if org.Domain == "" {
			http.Error(w, "Organization Domain Required", http.StatusBadRequest)
			return
		}

		query := `
INSERT INTO organisations (name, domain)
VALUES ($1, $2)
`

		_, err = db.Exec(query, org.Name, org.Domain)

		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"message": "Organization Created Successfully",
		}

		json.NewEncoder(w).Encode(response)

		return
	}

	if r.Method == http.MethodGet {

		json.NewEncoder(w).Encode(organisations)

		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
