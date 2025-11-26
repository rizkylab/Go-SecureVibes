package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
)

// Vulnerability 1: Hardcoded Credential
var dbPassword = "super_secret_password_123"

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/public/info", infoHandler)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	// Vulnerability 2: Weak Cryptography (MD5)
	h := md5.New()
	h.Write([]byte(username))
	fmt.Fprintf(w, "Hash: %x", h.Sum(nil))

	// Vulnerability 3: SQL Injection
	query := "SELECT * FROM users WHERE username = '" + username + "'"
	fmt.Println("Executing query:", query)

	// Mock DB execution
	// db.Exec(query)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	// Vulnerability 4: Sensitive Data Leak
	apiKey := "AKIAIOSFODNN7EXAMPLE"
	fmt.Fprintf(w, "System Info: API Key is %s", apiKey)
}
