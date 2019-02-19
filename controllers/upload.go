package controllers

import (
	"encoding/json"
	"net/http"

	"../nn"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type successMessage struct {
	Answer int `json:"answer"`
}

// GivePredict give predict from image file
func GivePredict(net *nn.NeuralNetwork) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		file, handle, err := r.FormFile("file")
		if err != nil {
			jsonErrorResponse(w, 1, "Empty file.")
			return
		}
		
		mimeType := handle.Header.Get("Content-Type")
		if mimeType != "image/png" {
			jsonErrorResponse(w, 2, "The format file is not valid.")
			return
		}

		answer := net.ImagePredict(file)
		jsonResponse(w, answer)
	}
}

func jsonResponse(w http.ResponseWriter, answer int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&successMessage{Answer: answer})
}

func jsonErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(&errorResponse{Message: message, Code: code})
}
