package goapidoc

// Document represents an api document information.
type Document struct {
	host     string
	basePath string
	info     *Info

	tags        []*Tag
	securities  []*Security
	paths       []*RoutePath
	definitions []*Definition
}

// NewDocument creates a default Document with given arguments.
func NewDocument(host, basePath string, info *Info) *Document {
	return &Document{host: host, basePath: basePath, info: info}
}

func (d *Document) GetHost() string               { return d.host }
func (d *Document) GetBasePath() string           { return d.basePath }
func (d *Document) GetInfo() *Info                { return d.info }
func (d *Document) GetTags() []*Tag               { return d.tags }
func (d *Document) GetSecurities() []*Security    { return d.securities }
func (d *Document) GetPaths() []*RoutePath        { return d.paths }
func (d *Document) GetDefinitions() []*Definition { return d.definitions }

// Tags sets the whole tags in Document.
func (d *Document) Tags(tags ...*Tag) *Document {
	d.tags = tags
	return d
}

// Tags adds some tags into Document.
func (d *Document) AddTags(tags ...*Tag) *Document {
	d.tags = append(d.tags, tags...)
	return d
}

// Securities sets the whole securities in Document.
func (d *Document) Securities(security ...*Security) *Document {
	d.securities = security
	return d
}

// Securities adds some securities into Document.
func (d *Document) AddSecurities(security ...*Security) *Document {
	d.securities = append(d.securities, security...)
	return d
}

// RoutePaths sets the whole route paths in Document.
func (d *Document) RoutePaths(paths ...*RoutePath) *Document {
	d.paths = paths
	return d
}

// AddRoutePaths adds some route paths into Document.
func (d *Document) AddRoutePaths(paths ...*RoutePath) *Document {
	d.paths = append(d.paths, paths...)
	return d
}

// Definitions sets the whole definitions in Document.
func (d *Document) Definitions(definitions ...*Definition) *Document {
	d.definitions = definitions
	return d
}

// AddDefinitions adds some definitions into Document.
func (d *Document) AddDefinitions(definitions ...*Definition) *Document {
	d.definitions = append(d.definitions, definitions...)
	return d
}

// Info represents a basic api information of Document.
type Info struct {
	title   string
	desc    string
	version string

	termsOfService string
	license        *License
	contact        *Contact
}

// NewInfo creates a default Info with given arguments.
func NewInfo(title, desc, version string) *Info {
	return &Info{title: title, desc: desc, version: version}
}

func (i *Info) GetTitle() string          { return i.title }
func (i *Info) GetDesc() string           { return i.desc }
func (i *Info) GetVersion() string        { return i.version }
func (i *Info) GetTermsOfService() string { return i.termsOfService }
func (i *Info) GetLicense() *License      { return i.license }
func (i *Info) GetContact() *Contact      { return i.contact }

// TermsOfService sets the termsOfService in Info.
func (i *Info) TermsOfService(service string) *Info {
	i.termsOfService = service
	return i
}

// License sets the license in Info.
func (i *Info) License(license *License) *Info {
	i.license = license
	return i
}

// Contact sets the contact in Info.
func (i *Info) Contact(contact *Contact) *Info {
	i.contact = contact
	return i
}

// License represents an api license information of Document.
type License struct {
	name string
	url  string
}

func (l *License) GetName() string { return l.name }
func (l *License) GetUrl() string  { return l.url }

// NewLicense creates a default License with given arguments.
func NewLicense(name, url string) *License {
	return &License{name: name, url: url}
}

// Contact represents an api contact information of Document.
type Contact struct {
	name  string
	url   string
	email string
}

func (c *Contact) GetName() string  { return c.name }
func (c *Contact) GetUrl() string   { return c.url }
func (c *Contact) GetEmail() string { return c.email }

// NewContact creates a default Contact with given arguments.
func NewContact(name, url, email string) *Contact {
	return &Contact{name: name, url: url, email: email}
}

// Tag represents an api tag information of Document.
type Tag struct {
	name string
	desc string
}

func (t *Tag) GetName() string { return t.name }
func (t *Tag) GetDesc() string { return t.desc }

// NewTag creates a default Tag with given arguments.
func NewTag(name, desc string) *Tag {
	return &Tag{name: name, desc: desc}
}

// Security represents an api security information of Document.
type Security struct {
	title string
	typ   string // only supports apiKey
	in    string
	name  string
}

func (s *Security) GetTitle() string { return s.title }
func (s *Security) GetIn() string    { return s.in }
func (s *Security) GetName() string  { return s.name }

// NewSecurity creates a default Security with given arguments.
func NewSecurity(title, in, name string) *Security {
	return &Security{title: title, typ: "apiKey", in: in, name: name}
}

// _document is the global Document.
var _document = NewDocument("", "", nil)

// Tags sets the whole tags in Document.
func SetTags(tags ...*Tag) *Document {
	return _document.Tags(tags...)
}

// Tags adds some tags into Document.
func AddTags(tags ...*Tag) *Document {
	return _document.AddTags(tags...)
}

// Securities sets the whole securities in Document.
func SetSecurities(security ...*Security) *Document {
	return _document.Securities(security...)
}

// Securities adds some securities into Document.
func AddSecurities(security ...*Security) *Document {
	return _document.AddSecurities(security...)
}

// RoutePaths sets the whole route paths in Document.
func SetRoutePaths(paths ...*RoutePath) *Document {
	return _document.RoutePaths(paths...)
}

// AddRoutePaths adds some route paths into Document.
func AddRoutePaths(paths ...*RoutePath) *Document {
	return _document.AddRoutePaths(paths...)
}

// Definitions sets the whole definitions in Document.
func SetDefinitions(definitions ...*Definition) *Document {
	return _document.Definitions(definitions...)
}

// AddDefinitions adds some definitions into Document.
func AddDefinitions(definitions ...*Definition) *Document {
	return _document.AddDefinitions(definitions...)
}
