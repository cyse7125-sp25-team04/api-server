package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webapp/services/constants"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var newUser User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received user: %+v\n", newUser)
	if newUser.FirstName == "" || newUser.LastName == "" || newUser.Username == "" || newUser.Password == "" || (newUser.Role != constants.UserRole && newUser.Role != constants.InstructorRole) {
		http.Error(w, "Missing required fields. Error: The following fileds are mandatory : firstName, lastName, username, passwordHash, role", http.StatusBadRequest)
		return
	}
	if IsUserExists(newUser.Username) {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	err = CreateUser(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// getUserHandler authenticates via basic auth and returns user details.
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authUser, err := Authenticate(username, password)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authUser)
}

// updateUserHandler authenticates the user and then updates the user's data.
// Returns 400 Bad Request for invalid payload/path, 401 for unauthorized,
// and 204 No Content on success.
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Verify basic auth credentials.
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authUser, err := Authenticate(username, password)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ensure the request has the expected Content-Type.
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// Decode the JSON payload into a User struct.
	var updatedUser User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if updatedUser.Username != username {
		http.Error(w, "can not update username", http.StatusBadRequest)
		return
	}
	// Force the update to use the authenticated user's ID.
	updatedUser.ID = authUser.ID

	// Update the user record.
	err = UpdateUser(&updatedUser)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// On success, return 204 No Content.
	w.WriteHeader(http.StatusNoContent)
}
