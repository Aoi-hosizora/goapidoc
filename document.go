package goapidoc

// Api document
type Document struct {
	host     string
	basePath string
	info     *Info

	tags        []*Tag
	securities  []*Security
	paths       []*RoutePath
	definitions []*Definition
}

func (d *Document) Tags(tags ...*Tag) *Document {
	d.tags = tags
	return d
}

func (d *Document) Securities(security ...*Security) *Document {
	d.securities = security
	return d
}

func (d *Document) AddRoutePaths(path ...*RoutePath) *Document {
	d.paths = append(d.paths, path...)
	return d
}

func (d *Document) AddDefinitions(models ...*Definition) *Document {
	d.definitions = append(d.definitions, models...)
	return d
}

// Base information
type Info struct {
	title   string
	desc    string
	version string

	termsOfService string
	license        *License
	contact        *Contact
}

func NewInfo(title string, desc string, version string) *Info {
	return &Info{title: title, desc: desc, version: version}
}

func (i *Info) TermsOfService(service string) *Info {
	i.termsOfService = service
	return i
}

func (i *Info) License(license *License) *Info {
	i.license = license
	return i
}

func (i *Info) Contact(contact *Contact) *Info {
	i.contact = contact
	return i
}

// License information
type License struct {
	name string
	url  string
}

func NewLicense(name string, url string) *License {
	return &License{name: name, url: url}
}

// Contact information
type Contact struct {
	name  string
	url   string
	email string
}

func NewContact(name string, url string, email string) *Contact {
	return &Contact{name: name, url: url, email: email}
}

// Security information
type Security struct {
	title string
	typ   string // only support for apiKey
	in    string
	name  string
}

func NewSecurity(title string, in string, name string) *Security {
	return &Security{title: title, typ: "apiKey", in: in, name: name}
}

// Tag information
type Tag struct {
	name string
	desc string
}

func NewTag(name string, desc string) *Tag {
	return &Tag{name: name, desc: desc}
}

// Global document
var _document = &Document{}

func SetDocument(host string, basePath string, info *Info) {
	_document.host = host
	_document.basePath = basePath
	_document.info = info
}

func SetTags(tags ...*Tag) {
	_document.Tags(tags...)
}

func SetSecurities(securities ...*Security) {
	_document.Securities(securities...)
}

func AddRoutePaths(paths ...*RoutePath) {
	_document.AddRoutePaths(paths...)
}

func AddDefinitions(definitions ...*Definition) {
	_document.AddDefinitions(definitions...)
}
