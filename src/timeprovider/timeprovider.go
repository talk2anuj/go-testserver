package main

import (
	"fmt"
	"net/http"
	"time"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>goto /gettime to get current time</h1>")
}

func gettime(w http.ResponseWriter, r *http.Request) {
	curr := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(w, curr)
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Health check</h1>")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/health_check", check)
	http.HandleFunc("/gettime", gettime)
	fmt.Println("Server starting...")
	http.ListenAndServe(":3000", nil)

	fmt.Println("here")
}
