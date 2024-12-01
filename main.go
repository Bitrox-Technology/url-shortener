package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // Convert the originalURl to the byte slice
	fmt.Println("Hasher: ", hasher)

	data := hasher.Sum(nil)
	fmt.Println("Data: ", data)

	hash := hex.EncodeToString(data)
	fmt.Println("hash: ", hash)
	fmt.Println("hash string: ", hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)

	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}

	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Method")

	// fmt.Fprintf(w, "Hello world")

}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	shortURL := createURL(data.URL)
	// fmt.Fprintf(w, shortURL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	fmt.Println("hash URl: ", id)
	url, err := getURL(id)
	fmt.Println("URL data: ", url)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {
	// fmt.Println("Starting URL Shortener...")
	// OriginalURL := "https://bitrox.tech/"
	// generateShortURL(OriginalURL)

	http.HandleFunc("/", handler)
	http.HandleFunc("/shortener", ShortURLHandler)
	http.HandleFunc("/redirect/{id}", RedirectURLHandler)
	// Start server
	fmt.Println("Starting Server On Port: 8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error on starting Server: ", err)
	}
}
