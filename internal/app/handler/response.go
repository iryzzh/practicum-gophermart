package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// failResponse формирует JSON ответ с ошибкой и кодом 400
func failResponse(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusBadRequest)

	type errResponse struct {
		Error string `json:"error"`
	}

	if e != nil {
		err := json.NewEncoder(w).Encode(errResponse{Error: e.Error()})
		if err != nil {
			log.Panic(err)
		}
	}
}

// basicResponse формирует 'plain-text' ответ с заданным кодом
func basicResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}
