package main

type LogStatus int64

const (
	EveryMinute           = "0 */1 * * * ?"
	Undo        LogStatus = -1
	Do          LogStatus = 1
)
