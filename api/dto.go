package main

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}
