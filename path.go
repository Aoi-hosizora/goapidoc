package yamldoc

type Path struct {
	Method string
	Route  string

	Summary     string
	Description string
	Tags        []string
	Consumes    []string
	Produces    []string
	Securities  []string
	Params      []*Param
	Responses   []*Response
	Deprecated  bool
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

func (r *Path) SetParams(params ...*Param) *Path {
	r.Params = params
	return r
}

func (r *Path) SetResponses(responses ...*Response) *Path {
	r.Responses = responses
	return r
}

func (r *Path) SetDeprecated(deprecated bool) *Path {
	r.Deprecated = deprecated
	return r
}

type Response struct {
	Code        int
	Description string
	Schema      string
	Headers     []*Header
	Examples    map[string]string
}

func NewResponse(code int) *Response {
	return &Response{Code: code}
}

func (r *Response) SetDescription(description string) *Response {
	r.Description = description
	return r
}

func (r *Response) SetSchema(ref string) *Response {
	r.Schema = ref
	return r
}

func (r *Response) SetHeaders(headers ...*Header) *Response {
	r.Headers = headers
	return r
}

func (r *Response) SetExamples(examples map[string]string) *Response {
	r.Examples = examples
	return r
}

type Header struct {
	Name        string
	Description string
	Type        string
	Format      string
	Default     interface{}
}

func NewHeader(name string, description string, headerType string) *Header {
	return &Header{Name: name, Description: description, Type: headerType}
}

func (h *Header) SetFormat(format string) *Header {
	h.Format = format
	return h
}

func (h *Header) SetDefault(defaultValue string) *Header {
	h.Default = defaultValue
	return h
}
