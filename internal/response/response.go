package response

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(statusCode int, message string) *Response {
	return &Response{
		Status:  statusCode,
		Message: message,
	}
}

func ErrorResponse(statusCode int, message string) *Response {
	return &Response{
		Status:  statusCode,
		Message: message,
	}
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}
