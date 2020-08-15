package goapidoc

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
	Code     int
	Type     string
	Required bool

	Description string
	Examples    map[string]string // content-type: example
	Headers     []*Header
}

func NewResponse(code int) *Response {
	return &Response{Code: code}
}

func (r *Response) WithType(t string) *Response {
	r.Type = t
	return r
}

func (r *Response) WithRequired(req bool) *Response {
	r.Required = req
	return r
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

// Response header
type Header struct {
	Name        string
	Type        string // base type
	Description string
}

func NewHeader(name string, t string, desc string) *Header {
	return &Header{Name: name, Type: t, Description: desc}
}

// Request parameter
type Param struct {
	Name        string
	In          string
	Type        string // string number integer boolean array (file)
	Required    bool
	Description string

	// `in` != body
	AllowEmptyValue bool
	Default         interface{}
	Example         interface{}
	Enum            []interface{}
	MinLength       int
	MaxLength       int
	Minimum         int
	Maximum         int
}

func NewParam(name string, in string, t string, req bool, desc string) *Param {
	return &Param{Name: name, In: in, Required: req, Description: desc, Type: t}
}

func NewPathParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, PATH, t, req, desc)
}

func NewQueryParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, QUERY, t, req, desc)
}

func NewFormParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, FORM, t, req, desc)
}

func NewBodyParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, BODY, t, req, desc)
}

func NewHeaderParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, HEADER, t, req, desc)
}

func (p *Param) WithAllowEmptyValue(allow bool) *Param {
	p.AllowEmptyValue = allow
	return p
}

func (p *Param) WithDefault(def interface{}) *Param {
	p.Default = def
	return p
}

func (p *Param) WithExample(ex interface{}) *Param {
	p.Example = ex
	return p
}

func (p *Param) WithEnum(enum ...interface{}) *Param {
	p.Enum = enum
	return p
}

func (p *Param) WithMinLength(min int) *Param {
	p.MinLength = min
	return p
}

func (p *Param) WithMaxLength(max int) *Param {
	p.MaxLength = max
	return p
}

func (p *Param) WithMinimum(min int) *Param {
	p.Minimum = min
	return p
}

func (p *Param) WithMaximum(max int) *Param {
	p.Maximum = max
	return p
}
