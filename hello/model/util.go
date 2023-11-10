package model

type ResponseResult struct {
	Code    int64
	Success bool
	Message string
	Data    interface{}
	Debug   string
}

type RequestParam struct {
	Msg string
}
