package yamldoc

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

func (r *Path) SetDescription(description string) *Path {
	r.Description = description
	return r
}

func (r *Path) SetTags(tags ...string) *Path {
	r.Tags = tags
	return r
}

func (r *Path) SetConsumes(consumes ...string) *Path {
	r.Consumes = consumes
	return r
}

func (r *Path) SetProduces(produces ...string) *Path {
	r.Produces = produces
	return r
}

func (r *Path) SetSecurities(securities ...string) *Path {
	r.Securities = securities
	return r
}

func (r *Path) SetDeprecated(deprecated bool) *Path {
	r.Deprecated = deprecated
	return r
}

// Set parameters
func (r *Path) SetParams(params ...*Param) *Path {
	r.Params = params
	return r
}

// Set responses
func (r *Path) SetResponses(responses ...*Response) *Path {
	r.Responses = responses
	return r
}

// Route response
type Response struct {
	Code int

	Description string
	Examples    map[string]string
	Headers     []*Header
	Schema      *Schema
}

func NewResponse(code int) *Response {
	return &Response{Code: code}
}

func (r *Response) SetDescription(description string) *Response {
	r.Description = description
	return r
}

func (r *Response) SetExamples(examples map[string]string) *Response {
	r.Examples = examples
	return r
}

func (r *Response) SetHeaders(headers ...*Header) *Response {
	r.Headers = headers
	return r
}

// !! Set schema, support objects, primitives and arrays
func (r *Response) SetSchema(schema *Schema) *Response {
	r.Schema = schema
	return r
}

// Response header
type Header struct {
	Name        string
	Type        string
	Description string
}

func NewHeader(name string, headerType string, description string) *Header {
	return &Header{Name: name, Description: description, Type: headerType}
}
