package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	w.Write([]byte("hello from login"))

	h.userService.Login()
}

func (h *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	fmt.Println("registration")
	w.Write([]byte("hello from registration"))

	h.userService.Registration()
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get users")
	w.Write([]byte("hello from get users"))

	h.userService.GetAllUsers()
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get user")
	w.Write([]byte("hello from get user"))

	h.userService.GetUser()
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update user")
	w.Write([]byte("hello from update user"))

	h.userService.UpdateUser()
}
