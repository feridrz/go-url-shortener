package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mu       sync.Mutex
)

func main() {
	http.HandleFunc("/", handlePost)
	http.HandleFunc("/redirect/", handleGet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handlePost handles the POST request for URL shortening
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	shortID := generateShortID()

	mu.Lock()
	urlStore[shortID] = originalURL
	mu.Unlock()

	shortenedURL := "http://localhost:8080/redirect/" + shortID
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(shortenedURL))
}

// handleGet handles the GET request for URL redirection
func handleGet(w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[len("/redirect/"):]

	mu.Lock()
	originalURL, exists := urlStore[shortID]
	mu.Unlock()

	if !exists {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// generateShortID generates a random string to use as a short ID
func generateShortID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
