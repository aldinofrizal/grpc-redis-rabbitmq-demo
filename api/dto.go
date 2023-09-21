package main

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type ErrorResponse struct {
	Detail string `json:"detail"`
}
