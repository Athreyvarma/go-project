package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID             int    `json:"id"`
	OrganisationID int    `json:"organisation_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

var users []User

func UserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		var user User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if user.OrganisationID == 0 {
			http.Error(w, "Organisation ID Required", http.StatusBadRequest)
			return
		}

		if user.Name == "" {
			http.Error(w, "User Name Required", http.StatusBadRequest)
			return
		}

		if user.Email == "" {
			http.Error(w, "User Email Required", http.StatusBadRequest)
			return
		}

		if user.Role == "" {
			http.Error(w, "User Role Required", http.StatusBadRequest)
			return
		}
		found := false

		for _, org := range organisations {

			if org.ID == user.OrganisationID {
				found = true
				break
			}

		}

		if !found {
			http.Error(w, "Organisation Not Found", http.StatusBadRequest)
			return
		}

		user.ID = len(users) + 1

		users = append(users, user)

		response := map[string]string{
			"message": "User Created Successfully",
		}

		json.NewEncoder(w).Encode(response)

		return
	}

	if r.Method == http.MethodGet {

		json.NewEncoder(w).Encode(users)

		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
