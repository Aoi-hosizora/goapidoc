package yamldoc

type Response struct {
	Code        int
	Description string
	Schema      string
	Examples    map[string]string
}

func NewResponse(code int) *Response {
	return &Response{Code: code}
}

func (r *Response) SetDescription(description string) *Response {
	r.Description = description
	return r
}

func (r *Response) SetSchema(schema string) *Response {
	r.Schema = schema
	return r
}

func (r *Response) SetExamples(examples map[string]string) *Response {
	r.Examples = examples
	return r
}
