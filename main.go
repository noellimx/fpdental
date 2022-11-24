package main

import http "net/http"

func main() {

	http.Handle("hi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	http.HandleFunc("a", func(w http.ResponseWriter, r *http.Request) {

	})
}
