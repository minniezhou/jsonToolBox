package jsonToolBox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1048576
	reader := http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&data)
	fmt.Println(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("should have only one json block")
	}
	return nil
}

func WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	jData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	_, err = w.Write(jData)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJson(w http.ResponseWriter, message string, statusCode ...int) error {
	status := http.StatusBadRequest
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	response := jsonResponse{
		Error:   true,
		Message: message,
	}
	return WriteJson(w, status, response)
}
