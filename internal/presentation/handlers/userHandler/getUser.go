package userhandler

import (
	"fmt"
	"net/http"

	"github.com/ciscapello/chat-backend/internal/presentation/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		response.SendError(w, http.StatusBadRequest, "Missing user ID parameter")
		uh.logErrorInRequest(r, "Missing user ID parameter")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid id")
		uh.logErrorInRequest(r, "invalid id")
		return
	}

	user, err := uh.userService.GetUser(uid)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "cannot get user")
		uh.logErrorInRequest(r, "cannot get user")
		return
	}

	response.SendSuccess(w, http.StatusOK, user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update user")
	w.Write([]byte("hello from update user"))

	uh.userService.UpdateUser()
}
