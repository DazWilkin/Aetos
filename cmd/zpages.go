package main

import (
	"net/http"
)

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	// Not ideal but there's no logger to use to handle the error
	_, _ = w.Write([]byte("ok"))
}
func statusz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	// Not ideal but there's no logger to use to handle the error
	_, _ = w.Write([]byte("ok"))
}
