package main

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInWeb()