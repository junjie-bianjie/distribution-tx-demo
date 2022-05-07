package main

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type MsgLog struct {
	UUID     string    `json:"uuid"`
	UserId   int64     `json:"user_id"`
	Integral int64     `json:"integral"`
	Status   LogStatus `json:"status"` // -1 没做，1做了
}

type Integral struct {
	Id       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	Integral int64 `json:"integral"`
}
