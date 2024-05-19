package userhandler

import (
	"encoding/json"
	"io"
	"net/http"
)

type requestBody struct {
	Username string
	Email    string
}

func (uh *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

	var rb requestBody
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		uh.logErrorInRequest(r, "Unable to read request body")
		return
	}

	if err := json.Unmarshal(body, &rb); err != nil {
		http.Error(w, "Unable to unmarshal request body", http.StatusBadRequest)
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
		return
	}

	err = uh.userService.Registration(rb.Username, rb.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		uh.logErrorInRequest(r, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
}
