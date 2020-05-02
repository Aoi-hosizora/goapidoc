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

type Param struct {
	Name        string
	In          string
	Type        string
	Required    bool
	Description string

	Format  string
	Default interface{}
	Schema  string
	Enum    []interface{}
}

func NewParam(name string, in string, paramType string, required bool, description string) *Param {
	return &Param{Name: name, In: in, Type: paramType, Required: required, Description: description}
}

func (p *Param) SetFormat(format string) *Param {
	p.Format = format
	return p
}

func (p *Param) SetDefault(defaultValue interface{}) *Param {
	p.Default = defaultValue
	return p
}

func (p *Param) SetSchema(schema string) *Param {
	p.Schema = schema
	return p
}

func (p *Param) SetEnum(enum ...interface{}) *Param {
	p.Enum = enum
	return p
}

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
