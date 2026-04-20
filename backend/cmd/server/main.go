package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}
