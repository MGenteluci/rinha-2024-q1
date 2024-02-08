package httputils

type HttpResponse struct {
	Status  int
	Headers map[string]string
	Body    interface{}
}

func NewHttpResponse(status int, body interface{}) *HttpResponse {
	var b interface{}
	if body != nil {
		b = body
	}
	return &HttpResponse{
		Status: status,
		Body: b,
	}
}

func Ok(body interface{}) *HttpResponse {
	return NewHttpResponse(200, body)
}

func NotFound() *HttpResponse {
	return NewHttpResponse(404, nil)
}

func MethodNotAllowed() *HttpResponse {
	return NewHttpResponse(405, nil)
}

func UnprocessableEntity(body interface{}) *HttpResponse {
	return NewHttpResponse(422, body)
}

func InternalServerError(body interface{}) *HttpResponse {
	return NewHttpResponse(500, body)
}
