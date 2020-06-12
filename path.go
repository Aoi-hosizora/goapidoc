package apidoc

// Route path
type Path struct {
	Method  string
	Route   string
	Summary string

	Description string
	Tags        []string
	Consumes    []string
	Produces    []string
	Securities  []string
	Deprecated  bool

	Params    []*Param
	Responses []*Response
}

func NewPath(method string, route string, summary string) *Path {
	return &Path{Method: method, Route: route, Summary: summary}
}

func (r *Path) WithDescription(desc string) *Path {
	r.Description = desc
	return r
}

func (r *Path) WithTags(tags ...string) *Path {
	r.Tags = tags
	return r
}

func (r *Path) WithConsumes(consumes ...string) *Path {
	r.Consumes = consumes
	return r
}

func (r *Path) WithProduces(produces ...string) *Path {
	r.Produces = produces
	return r
}

func (r *Path) WithSecurities(securities ...string) *Path {
	r.Securities = securities
	return r
}

func (r *Path) WithDeprecated(deprecated bool) *Path {
	r.Deprecated = deprecated
	return r
}

// Set parameters
func (r *Path) WithParams(params ...*Param) *Path {
	r.Params = params
	return r
}

// Set responses
func (r *Path) WithResponses(responses ...*Response) *Path {
	r.Responses = responses
	return r
}

// Route response
type Response struct {
	Code int

	Description string
	Examples    map[string]string // content-type: example
	Headers     []*Header
	Schema      *Schema
}

func NewResponse(code int) *Response {
	return &Response{Code: code}
}

func (r *Response) WithDescription(desc string) *Response {
	r.Description = desc
	return r
}

func (r *Response) WithExamples(examples map[string]string) *Response {
	r.Examples = examples
	return r
}

func (r *Response) WithHeaders(headers ...*Header) *Response {
	r.Headers = headers
	return r
}

// !! Set schema, support objects, primitives and arrays
// Deprecated
func (r *Response) WithSchema(schema *Schema) *Response {
	r.Schema = schema
	return r
}

// Response header
type Header struct {
	Name        string
	Type        string
	Description string
}

func NewHeader(name string, t string, desc string) *Header {
	return &Header{Name: name, Type: t, Description: desc}
}
