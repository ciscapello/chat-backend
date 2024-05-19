package userhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get user")
	w.Write([]byte("hello from get user"))

	var id string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		uh.logErrorInRequest(r, "Unable to read request body")
	}

	if err := json.Unmarshal(body, &id); err != nil {
		http.Error(w, "Unable to unmarshal request body", http.StatusBadRequest)
		uh.logErrorInRequest(r, "Unable to unmarshal request body")
	}

	uuid, err := uuid.Parse(string(id))
	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusBadRequest)
		uh.logErrorInRequest(r, "Unable to parse request body")
	}

	uh.userService.GetUser(uuid)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update user")
	w.Write([]byte("hello from update user"))

	uh.userService.UpdateUser()
}
