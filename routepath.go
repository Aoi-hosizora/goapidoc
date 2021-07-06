package goapidoc

// =========
// RoutePath
// =========

// RoutePath represents an api route path information of Document.
type RoutePath struct {
	method  string
	route   string
	summary string

	desc        string
	operationId string
	schemas     []string
	consumes    []string
	produces    []string
	tags        []string
	securities  []string
	deprecated  bool
	params      []*Param
	responses   []*Response
}

// NewRoutePath creates a default RoutePath with given arguments.
func NewRoutePath(method, route, summary string) *RoutePath {
	return &RoutePath{method: method, route: route, summary: summary}
}

// NewGetRoutePath creates a get RoutePath with given arguments.
func NewGetRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(GET, route, summary)
}

// NewPutRoutePath creates a put RoutePath with given arguments.
func NewPutRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(PUT, route, summary)
}

// NewPostRoutePath creates a post RoutePath with given arguments.
func NewPostRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(POST, route, summary)
}

// NewDeleteRoutePath creates a delete RoutePath with given arguments.
func NewDeleteRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(DELETE, route, summary)
}

// NewOptionsRoutePath creates a options RoutePath with given arguments.
func NewOptionsRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(OPTIONS, route, summary)
}

// NewHeadRoutePath creates a head RoutePath with given arguments.
func NewHeadRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(HEAD, route, summary)
}

// NewPatchRoutePath creates a patch RoutePath with given arguments.
func NewPatchRoutePath(route, summary string) *RoutePath {
	return NewRoutePath(PATCH, route, summary)
}

func (r *RoutePath) GetMethod() string         { return r.method }
func (r *RoutePath) GetRoute() string          { return r.route }
func (r *RoutePath) GetSummary() string        { return r.summary }
func (r *RoutePath) GetDesc() string           { return r.desc }
func (r *RoutePath) GetOperationId() string    { return r.operationId }
func (r *RoutePath) GetSchemas() []string      { return r.schemas }
func (r *RoutePath) GetConsumes() []string     { return r.consumes }
func (r *RoutePath) GetProduces() []string     { return r.produces }
func (r *RoutePath) GetTags() []string         { return r.tags }
func (r *RoutePath) GetSecurities() []string   { return r.securities }
func (r *RoutePath) GetDeprecated() bool       { return r.deprecated }
func (r *RoutePath) GetParams() []*Param       { return r.params }
func (r *RoutePath) GetResponses() []*Response { return r.responses }

// Method sets the method in RoutePath.
func (r *RoutePath) Method(method string) *RoutePath {
	r.method = method
	return r
}

// Route sets the route in RoutePath.
func (r *RoutePath) Route(route string) *RoutePath {
	r.route = route
	return r
}

// Summary sets the summary in RoutePath.
func (r *RoutePath) Summary(summary string) *RoutePath {
	r.summary = summary
	return r
}

// Desc sets the desc in RoutePath.
func (r *RoutePath) Desc(desc string) *RoutePath {
	r.desc = desc
	return r
}

// OperationId sets the operationId in RoutePath.
func (r *RoutePath) OperationId(operationId string) *RoutePath {
	r.operationId = operationId
	return r
}

// Schemas sets the whole schemas in RoutePath.
func (r *RoutePath) Schemas(schemas ...string) *RoutePath {
	r.schemas = schemas
	return r
}

// AddSchemas adds some tags schemas into RoutePath.
func (r *RoutePath) AddSchemas(schemas ...string) *RoutePath {
	r.schemas = append(r.schemas, schemas...)
	return r
}

// Consumes sets the whole consumes in RoutePath.
func (r *RoutePath) Consumes(consumes ...string) *RoutePath {
	r.consumes = consumes
	return r
}

// AddConsumes adds some consumes into RoutePath.
func (r *RoutePath) AddConsumes(consumes ...string) *RoutePath {
	r.consumes = append(r.consumes, consumes...)
	return r
}

// Produces sets the whole produces in RoutePath.
func (r *RoutePath) Produces(produces ...string) *RoutePath {
	r.produces = produces
	return r
}

// AddProduces adds some produces into RoutePath.
func (r *RoutePath) AddProduces(produces ...string) *RoutePath {
	r.produces = append(r.produces, produces...)
	return r
}

// Tags sets the whole tags in RoutePath.
func (r *RoutePath) Tags(tags ...string) *RoutePath {
	r.tags = tags
	return r
}

// AddTags adds some tags into RoutePath.
func (r *RoutePath) AddTags(tags ...string) *RoutePath {
	r.tags = append(r.tags, tags...)
	return r
}

// Securities sets the whole security-requirements in RoutePath.
func (r *RoutePath) Securities(securities ...string) *RoutePath {
	r.securities = securities
	return r
}

// AddSecurities adds some security-requirements into RoutePath.
func (r *RoutePath) AddSecurities(securities ...string) *RoutePath {
	r.securities = append(r.securities, securities...)
	return r
}

// Deprecated sets the deprecated in RoutePath.
func (r *RoutePath) Deprecated(deprecated bool) *RoutePath {
	r.deprecated = deprecated
	return r
}

// Params sets the whole params in RoutePath.
func (r *RoutePath) Params(params ...*Param) *RoutePath {
	r.params = params
	return r
}

// AddParams adds some params into RoutePath.
func (r *RoutePath) AddParams(params ...*Param) *RoutePath {
	r.params = append(r.params, params...)
	return r
}

// Responses sets the whole responses in RoutePath.
func (r *RoutePath) Responses(responses ...*Response) *RoutePath {
	r.responses = responses
	return r
}

// AddResponses adds some responses into RoutePath.
func (r *RoutePath) AddResponses(responses ...*Response) *RoutePath {
	r.responses = append(r.responses, responses...)
	return r
}

// ========
// Response
// ========

// Response represents a route response information of RoutePath.
type Response struct {
	code int
	typ  string

	desc     string
	examples map[string]string
	headers  []*Header
}

// NewResponse creates a default Response with given arguments.
func NewResponse(code int, typ string) *Response {
	return &Response{code: code, typ: typ}
}

func (r *Response) GetCode() int                   { return r.code }
func (r *Response) GetType() string                { return r.typ }
func (r *Response) GetDesc() string                { return r.desc }
func (r *Response) GetExamples() map[string]string { return r.examples }
func (r *Response) GetHeaders() []*Header          { return r.headers }

// Code sets the code in Response.
func (r *Response) Code(code int) *Response {
	r.code = code
	return r
}

// Type sets the type in Response.
func (r *Response) Type(typ string) *Response {
	r.typ = typ
	return r
}

// Desc sets the desc in Response.
func (r *Response) Desc(desc string) *Response {
	r.desc = desc
	return r
}

// Examples sets the whole examples in Response.
func (r *Response) Examples(examples map[string]string) *Response {
	r.examples = examples
	return r
}

// Headers sets the whole headers in Response.
func (r *Response) Headers(headers ...*Header) *Response {
	r.headers = headers
	return r
}

// AddHeaders add some headers in Response.
func (r *Response) AddHeaders(headers ...*Header) *Response {
	r.headers = append(r.headers, headers...)
	return r
}

// ======
// Header
// ======

// Header represents a response header information of Response.
type Header struct {
	name string
	typ  string // primitive type
	desc string
}

// NewHeader creates a default Header with given arguments.
func NewHeader(name, typ, desc string) *Header {
	return &Header{name: name, typ: typ, desc: desc}
}

func (h *Header) GetName() string { return h.name }
func (h *Header) GetType() string { return h.typ }
func (h *Header) GetDesc() string { return h.desc }

// Name sets the name in Header.
func (h *Header) Name(name string) *Header {
	h.name = name
	return h
}

// Type sets the type in Header.
func (h *Header) Type(typ string) *Header {
	h.typ = typ
	return h
}

// Desc sets the desc in Header.
func (h *Header) Desc(desc string) *Header {
	h.desc = desc
	return h
}

// =====
// Param
// =====

// Param represents a request parameter information of RoutePath.
type Param struct {
	name     string
	in       string
	typ      string // string number integer boolean array file
	required bool
	desc     string

	allowEmpty       bool
	defaul           interface{}
	example          interface{}
	pattern          string
	enum             []interface{}
	minLength        int
	maxLength        int
	minItems         int
	maxItems         int
	uniqueItems      bool
	collectionFormat string
	minimum          float64
	maximum          float64
	exclusiveMin     bool
	exclusiveMax     bool
	multipleOf       float64
	itemOption       *ItemOption
}

// NewParam creates a default Param with given arguments.
func NewParam(name, in, typ string, required bool, desc string) *Param {
	return &Param{name: name, in: in, typ: typ, required: required, desc: desc}
}

// NewPathParam creates a path Param with given arguments.
func NewPathParam(name, typ string, required bool, desc string) *Param {
	return NewParam(name, PATH, typ, required, desc)
}

// NewQueryParam creates a query Param with given arguments.
func NewQueryParam(name, typ string, required bool, desc string) *Param {
	return NewParam(name, QUERY, typ, required, desc)
}

// NewFormParam creates a form Param with given arguments.
func NewFormParam(name, typ string, required bool, desc string) *Param {
	return NewParam(name, FORM, typ, required, desc)
}

// NewBodyParam creates a body Param with given arguments.
func NewBodyParam(name, typ string, required bool, desc string) *Param {
	return NewParam(name, BODY, typ, required, desc)
}

// NewHeaderParam creates a header Param with given arguments.
func NewHeaderParam(name, typ string, required bool, desc string) *Param {
	return NewParam(name, HEADER, typ, required, desc)
}

func (p *Param) GetName() string             { return p.name }
func (p *Param) GetInLoc() string            { return p.in }
func (p *Param) GetType() string             { return p.typ }
func (p *Param) GetRequired() bool           { return p.required }
func (p *Param) GetAllowEmpty() bool         { return p.allowEmpty }
func (p *Param) GetDefault() interface{}     { return p.defaul }
func (p *Param) GetExample() interface{}     { return p.example }
func (p *Param) GetPattern() string          { return p.pattern }
func (p *Param) GetEnum() []interface{}      { return p.enum }
func (p *Param) GetMinLength() int           { return p.minLength }
func (p *Param) GetMaxLength() int           { return p.maxLength }
func (p *Param) GetMinItems() int            { return p.minItems }
func (p *Param) GetMaxItems() int            { return p.maxItems }
func (p *Param) GetUniqueItems() bool        { return p.uniqueItems }
func (p *Param) GetCollectionFormat() string { return p.collectionFormat }
func (p *Param) GetMinimum() float64         { return p.minimum }
func (p *Param) GetMaximum() float64         { return p.maximum }
func (p *Param) GetExclusiveMin() bool       { return p.exclusiveMin }
func (p *Param) GetExclusiveMax() bool       { return p.exclusiveMax }
func (p *Param) GetMultipleOf() float64      { return p.multipleOf }
func (p *Param) GetItemOption() *ItemOption  { return p.itemOption }

// Name sets the name in Param.
func (p *Param) Name(name string) *Param {
	p.name = name
	return p
}

// InLoc sets the in-location in Param.
func (p *Param) InLoc(in string) *Param {
	p.in = in
	return p
}

// Type sets the type in Param.
func (p *Param) Type(typ string) *Param {
	p.typ = typ
	return p
}

// Required sets the required in Param.
func (p *Param) Required(required bool) *Param {
	p.required = required
	return p
}

// Desc sets the desc in Param.
func (p *Param) Desc(desc string) *Param {
	p.desc = desc
	return p
}

// AllowEmpty sets the allowEmpty in Param.
func (p *Param) AllowEmpty(allow bool) *Param {
	p.allowEmpty = allow
	return p
}

// Default sets the default in Param.
func (p *Param) Default(defaul interface{}) *Param {
	p.defaul = defaul
	return p
}

// Example sets the example in Param.
func (p *Param) Example(example interface{}) *Param {
	p.example = example
	return p
}

// Pattern sets the pattern in Param.
func (p *Param) Pattern(pattern string) *Param {
	p.pattern = pattern
	return p
}

// Enum sets the whole enums in Param.
func (p *Param) Enum(enums ...interface{}) *Param {
	p.enum = enums
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

// LengthRange sets the minLength and maxLength in Param.
// TODO BREAK CHANGES
func (p *Param) LengthRange(min, max int) *Param {
	p.minLength = min
	p.maxLength = max
	return p
}

// MinItems sets the minItems in Param.
func (p *Param) MinItems(min int) *Param {
	p.minItems = min
	return p
}

// MaxItems sets the maxItems in Param.
func (p *Param) MaxItems(max int) *Param {
	p.maxItems = max
	return p
}

// ItemsRange sets the minItems and maxItems in Param.
func (p *Param) ItemsRange(min, max int) *Param {
	p.minItems = min
	p.maxItems = max
	return p
}

// UniqueItems sets the uniqueItems in Param.
func (p *Param) UniqueItems(unique bool) *Param {
	p.uniqueItems = unique
	return p
}

// CollectionFormat sets the collectionFormat in Param.
func (p *Param) CollectionFormat(collectionFormat string) *Param {
	p.collectionFormat = collectionFormat
	return p
}

// Minimum sets the minimum in Param.
// TODO BREAK CHANGES
func (p *Param) Minimum(min float64) *Param {
	p.minimum = min
	return p
}

// Maximum sets the maximum in Param.
// TODO BREAK CHANGES
func (p *Param) Maximum(max float64) *Param {
	p.maximum = max
	return p
}

// ValueRange sets the minimum and maximum in Param.
// TODO BREAK CHANGES
func (p *Param) ValueRange(min, max float64) *Param {
	p.minimum = min
	p.maximum = max
	return p
}

// ExclusiveMin sets the exclusiveMin in Param.
func (p *Param) ExclusiveMin(exclusiveMin bool) *Param {
	p.exclusiveMin = exclusiveMin
	return p
}

// ExclusiveMax sets the exclusiveMax in Param.
func (p *Param) ExclusiveMax(exclusiveMax bool) *Param {
	p.exclusiveMax = exclusiveMax
	return p
}

// MultipleOf sets the multipleOf in Param.
func (p *Param) MultipleOf(multipleOf float64) *Param {
	p.multipleOf = multipleOf
	return p
}

// ItemOption sets the itemOption in ItemOption.
func (p *Param) ItemOption(itemOption *ItemOption) *Param {
	p.itemOption = itemOption
	return p
}
