package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// create a user struct
type User struct {
	// Username and Password are the fields for the User struct, they are exported so they can be accessed outside the package for valid JSON encoding/decoding
	Username string `json:"username"`
	Password string `json:"password"`
}

// users is a map to store users, where the key is an integer ID and the value is a User struct
var users = make(map[int]User)

func main() {
	//mux to handle HTTP requests
	mux := http.NewServeMux()
	// Register the handler function for the root path, this is GET by default
	mux.HandleFunc("/", handleRoot)
	// Register the handler function for the POST request to /users
	mux.HandleFunc("POST /users", createUser)
	fmt.Println("Starting server on :8080")
	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", mux)
}

// w is the response writer, r is the request
// w contains the response to be sent back to the client
// r contains the request information from the client
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!\n")
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
	// Store the user in the map with an incremented key
	users[len(users)+1] = user
	//sets response to 204
	w.WriteHeader(http.StatusNoContent)

}
