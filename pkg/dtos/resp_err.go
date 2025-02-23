package dtos

type RespErr struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}
