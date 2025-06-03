package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// create a user struct
type User struct {
	// Username and Password are the fields for the User struct, they are exported so they can be accessed outside the package for valid JSON encoding/decoding
	Username string `json:"username"`
	Password string `json:"password"`
}

// users is a map to store users, where the key is an integer ID and the value is a User struct
var users = make(map[int]User)

// cacheMutex is a mutex to protect the users map from concurrent access
// It ensures that only one goroutine can access the users map at a time, preventing conflicts
var cacheMutex sync.Mutex

func main() {
	//mux to handle HTTP requests
	mux := http.NewServeMux()
	// Register the handler function for the root path, this is GET by default
	mux.HandleFunc("/", handleRoot)
	// Register the handler function for the POST request to /users
	mux.HandleFunc("POST /users", createUser)
	// Register the handler function for the GET request to /users/{id}, id is a path parameter
	mux.HandleFunc("GET /users/{id}", getUser)
	// Delete request to delete a user by ID
	mux.HandleFunc("DELETE /users/{id}", deleteUser)
	fmt.Println("Starting server on :8080")
	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", mux)
}

// w is the response writer, r is the request
// w contains the response to be sent back to the client
// r contains the request information from the client
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User Manager Home!\n")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	//decode the body of the request into the User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	// if there is an error decoding the request body, return a 400 Bad Request error
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Check if the username and password are not empty, if they are, return a 400 Bad Request error
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Lock the cacheMutex to ensure that only one goroutine can access the users map at a time
	cacheMutex.Lock()
	// Store the user in the map with an incremented key
	users[len(users)+1] = user
	// Unlock the cacheMutex to allow other goroutines to access the users map
	cacheMutex.Unlock()
	//sets response to 204
	w.WriteHeader(http.StatusNoContent)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the path parameter, which is expected to be an integer
	id, err := strconv.Atoi(r.PathValue("id"))

	// If there is an error converting the ID to an integer, return a 400 Bad Request error
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Lock the cacheMutex to ensure that only one goroutine can access the users map at a time
	cacheMutex.Lock()
	user, ok := users[id]
	cacheMutex.Unlock()

	// If the user with the given ID does not exist, return a 404 Not Found error
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// If the user exists, encode the user struct to JSON and write it to the response
	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
		return
	}
	//return 200
	w.WriteHeader(http.StatusOK)
	// Set the Content-Type header to application/json and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response to the response writer
	w.Write(j)

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the path parameter, which is expected to be an integer
	id, err := strconv.Atoi(r.PathValue("id"))

	// If there is an error converting the ID to an integer, return a 400 Bad Request error
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// check if id is in there
	if _, ok := users[id]; !ok {
		// If the user with the given ID does not exist, return a 404 Not Found error
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Lock the cacheMutex to ensure that only one goroutine can access the users map at a time
	cacheMutex.Lock()
	delete(users, id)   // Delete the user from the map
	cacheMutex.Unlock() // Unlock the cacheMutex to allow other goroutines to access the users map

	w.WriteHeader(http.StatusNoContent)
}
