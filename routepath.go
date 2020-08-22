package goapidoc

// Route path
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

func NewRoutePath(method string, route string, summary string) *RoutePath {
	return &RoutePath{method: method, route: route, summary: summary}
}

func (r *RoutePath) Desc(desc string) *RoutePath {
	r.desc = desc
	return r
}

func (r *RoutePath) Tags(tags ...string) *RoutePath {
	r.tags = tags
	return r
}

func (r *RoutePath) Consumes(consumes ...string) *RoutePath {
	r.consumes = consumes
	return r
}

func (r *RoutePath) Produces(produces ...string) *RoutePath {
	r.produces = produces
	return r
}

func (r *RoutePath) Securities(securities ...string) *RoutePath {
	r.securities = securities
	return r
}

func (r *RoutePath) Deprecated(deprecated bool) *RoutePath {
	r.deprecated = deprecated
	return r
}

func (r *RoutePath) Params(params ...*Param) *RoutePath {
	r.params = params
	return r
}

func (r *RoutePath) Responses(responses ...*Response) *RoutePath {
	r.responses = responses
	return r
}

// Route response
type Response struct {
	code     int
	typ      string
	required bool

	desc     string
	examples map[string]string // content-type: example
	headers  []*Header
}

func NewResponse(code int) *Response {
	return &Response{code: code}
}

func (r *Response) Type(typ string) *Response {
	r.typ = typ
	return r
}

func (r *Response) Required(req bool) *Response {
	r.required = req
	return r
}

func (r *Response) Desc(desc string) *Response {
	r.desc = desc
	return r
}

func (r *Response) Examples(examples map[string]string) *Response {
	r.examples = examples
	return r
}

func (r *Response) Headers(headers ...*Header) *Response {
	r.headers = headers
	return r
}

// Response header
type Header struct {
	name string
	typ  string // base type
	desc string
}

func NewHeader(name string, typ string, desc string) *Header {
	return &Header{name: name, typ: typ, desc: desc}
}

// Request parameter
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

func NewParam(name string, in string, typ string, req bool, desc string) *Param {
	return &Param{name: name, in: in, required: req, desc: desc, typ: typ}
}

func NewPathParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, PATH, typ, req, desc)
}

func NewQueryParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, QUERY, typ, req, desc)
}

func NewFormParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, FORM, typ, req, desc)
}

func NewBodyParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, BODY, typ, req, desc)
}

func NewHeaderParam(name string, typ string, req bool, desc string) *Param {
	return NewParam(name, HEADER, typ, req, desc)
}

func (p *Param) AllowEmpty(allow bool) *Param {
	p.allowEmpty = allow
	return p
}

func (p *Param) Default(def interface{}) *Param {
	p.def = def
	return p
}

func (p *Param) Example(ex interface{}) *Param {
	p.example = ex
	return p
}

func (p *Param) Enum(enum ...interface{}) *Param {
	p.enum = enum
	return p
}

func (p *Param) MinLength(min int) *Param {
	p.minLength = min
	return p
}

func (p *Param) MaxLength(max int) *Param {
	p.maxLength = max
	return p
}

func (p *Param) Minimum(min int) *Param {
	p.minimum = min
	return p
}

func (p *Param) Maximum(max int) *Param {
	p.maximum = max
	return p
}
