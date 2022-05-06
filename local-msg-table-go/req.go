package main

type UserReq struct {
	Username string `json:"username" binding:"required"`
}
