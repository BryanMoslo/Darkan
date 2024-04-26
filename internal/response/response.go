package response

type Response struct {
	Status   int         `json:"status"`
	Message  string      `json:"message,omitempty"`
	Messages []string    `json:"messages,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func SuccessResponse(statusCode int, message string) *Response {
	return &Response{
		Status:  statusCode,
		Message: message,
	}
}

func ErrorResponse(statusCode int, messages ...string) *Response {
	return &Response{
		Status:   statusCode,
		Messages: messages,
	}
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}
