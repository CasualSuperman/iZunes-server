package main

import (
	"errors"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) (Song, int, error) {
	if r.Method != "POST" {
		return Song{}, http.StatusMethodNotAllowed, errors.New("upload must be post")
	}
	return Song{}, http.StatusOK, nil
}
