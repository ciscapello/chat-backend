package userhandler

import "net/http"

func (uh *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

	uh.userService.Registration()
}
