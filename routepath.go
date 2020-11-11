package goapidoc

// RoutePath is a route path information.
type RoutePath struct {
	method  string
	route   string
	summary string

	desc       string
	tags       []string
	consumes   []string
	produces   []string
	securities []string
	deprecated bool

	params    []*Param
	responses []*Response
}

// NewRoutePath creates a RoutePath.
func NewRoutePath(method string, route string, summary string) *RoutePath {
	return &RoutePath{method: method, route: route, summary: summary}
}

// Desc sets the desc in RoutePath.
func (r *RoutePath) Desc(desc string) *RoutePath {
	r.desc = desc
	return r
}

// Tags sets the tags in RoutePath.
func (r *RoutePath) Tags(tags ...string) *RoutePath {
	r.tags = tags
	return r
}

// Consumes sets the consumes in RoutePath.
func (r *RoutePath) Consumes(consumes ...string) *RoutePath {
	r.consumes = consumes
	return r
}

// Produces sets the produces in RoutePath.
func (r *RoutePath) Produces(produces ...string) *RoutePath {
	r.produces = produces
	return r
}

// Securities sets the securities in RoutePath.
func (r *RoutePath) Securities(securities ...string) *RoutePath {
	r.securities = securities
	return r
}

// Deprecated sets the deprecated in RoutePath.
func (r *RoutePath) Deprecated(deprecated bool) *RoutePath {
	r.deprecated = deprecated
	return r
}

// Params sets the params in RoutePath.
func (r *RoutePath) Params(params ...*Param) *RoutePath {
	r.params = params
	return r
}

// Responses sets the responses in RoutePath.
func (r *RoutePath) Responses(responses ...*Response) *RoutePath {
	r.responses = responses
	return r
}

// Response is a response information.
type Response struct {
	code int
	typ  string

	desc     string
	examples map[string]string // content-type: example
	headers  []*Header
}

// NewResponse creates a Response information.
func NewResponse(code int, typ string) *Response {
	return &Response{code: code, typ: typ}
}

// Desc sets the desc in RoutePath.
func (r *Response) Desc(desc string) *Response {
	r.desc = desc
	return r
}

// Examples sets the examples in RoutePath.
func (r *Response) Examples(examples map[string]string) *Response {
	r.examples = examples
	return r
}

// Headers sets the headers in RoutePath.
func (r *Response) Headers(headers ...*Header) *Response {
	r.headers = headers
	return r
}

// Header is a response header information.
type Header struct {
	name string
	typ  string // base type
	desc string
}

// NewHeader creates a Header.
func NewHeader(name string, typ string, desc string) *Header {
	return &Header{name: name, typ: typ, desc: desc}
}

// Param is a request parameter information.
type Param struct {
	name     string
	in       string
	typ      string // string number integer boolean array (file)
	required bool
	desc     string

	// `in` != body
	allowEmpty bool
	def        interface{}
	example    interface{}
	enum       []interface{}
	minLength  int
	maxLength  int
	minimum    int
	maximum    int
}

// NewParam creates a Param.
func NewParam(name string, in string, typ string, req bool, desc string) *Param {
	return &Param{name: name, in: in, required: req, desc: desc, typ: typ}
}

// NewPathParam creates a Param with path in-type.
func NewPathParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, PATH, typ, req, desc)
}

// NewQueryParam creates a Param with query in-type.
func NewQueryParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, QUERY, typ, req, desc)
}

// NewFormParam creates a Param with form in-type.
func NewFormParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, FORM, typ, req, desc)
}

// NewBodyParam creates a Param with body in-type.
func NewBodyParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, BODY, typ, req, desc)
}

// NewHeaderParam creates a Param with head in-type.
func NewHeaderParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, HEADER, typ, req, desc)
}

// AllowEmpty sets the allowEmpty in Param.
func (p *Param) AllowEmpty(allow bool) *Param {
	p.allowEmpty = allow
	return p
}

// Default sets the default in Param.
func (p *Param) Default(def interface{}) *Param {
	p.def = def
	return p
}

// Example sets the example in Param.
func (p *Param) Example(ex interface{}) *Param {
	p.example = ex
	return p
}

// Enum sets the enum in Param.
func (p *Param) Enum(enum ...interface{}) *Param {
	p.enum = enum
	return p
}

// MinLength sets the minLength in Param.
func (p *Param) MinLength(min int) *Param {
	p.minLength = min
	return p
}

// MaxLength sets the maxLength in Param.
func (p *Param) MaxLength(max int) *Param {
	p.maxLength = max
	return p
}

// Length sets the minLength and maxLength in Param.
func (p *Param) Length(min, max int) *Param {
	p.minLength = min
	p.maxLength = max
	return p
}

// Minimum sets the minimum in Param.
func (p *Param) Minimum(min int) *Param {
	p.minimum = min
	return p
}

// Maximum sets the maximum in Param.
func (p *Param) Maximum(max int) *Param {
	p.maximum = max
	return p
}

// Minimum sets the minimum and maximum in Param.
func (p *Param) MinMaximum(min, max int) *Param {
	p.minimum = min
	p.maximum = max
	return p
}
