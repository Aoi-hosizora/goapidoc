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

func (r *Path) SetParam(params ...*Param) *Path {
	r.Params = params
	return r
}

func (r *Path) SetResponse(responses ...*Response) *Path {
	r.Responses = responses
	return r
}
