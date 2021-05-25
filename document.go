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

// Host sets the host in Document.
func (d *Document) Host(host string) *Document {
	d.host = host
	return d
}

// BasePath sets the basePath in Document.
func (d *Document) BasePath(basePath string) *Document {
	d.basePath = basePath
	return d
}

// Info sets the info in Document.
func (d *Document) Info(info *Info) *Document {
	d.info = info
	return d
}

// Tags sets the whole tags in Document.
func (d *Document) Tags(tags ...*Tag) *Document {
	d.tags = tags
	return d
}

// AddTags adds some tags into Document.
func (d *Document) AddTags(tags ...*Tag) *Document {
	d.tags = append(d.tags, tags...)
	return d
}

// Securities sets the whole securities in Document.
func (d *Document) Securities(securities ...*Security) *Document {
	d.securities = securities
	return d
}

// AddSecurities adds some securities into Document.
func (d *Document) AddSecurities(securities ...*Security) *Document {
	d.securities = append(d.securities, securities...)
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

// Title sets the title in Info.
func (i *Info) Title(title string) *Info {
	i.title = title
	return i
}

// Desc sets the desc in Info.
func (i *Info) Desc(desc string) *Info {
	i.desc = desc
	return i
}

// Version sets the version in Info.
func (i *Info) Version(version string) *Info {
	i.version = version
	return i
}

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

// NewLicense creates a default License with given arguments.
func NewLicense(name, url string) *License {
	return &License{name: name, url: url}
}

func (l *License) GetName() string { return l.name }
func (l *License) GetUrl() string  { return l.url }

// Name sets the name in License.
func (l *License) Name(name string) *License {
	l.name = name
	return l
}

// Url sets the url in License.
func (l *License) Url(url string) *License {
	l.url = url
	return l
}

// Contact represents an api contact information of Document.
type Contact struct {
	name  string
	url   string
	email string
}

// NewContact creates a default Contact with given arguments.
func NewContact(name, url, email string) *Contact {
	return &Contact{name: name, url: url, email: email}
}

func (c *Contact) GetName() string  { return c.name }
func (c *Contact) GetUrl() string   { return c.url }
func (c *Contact) GetEmail() string { return c.email }

// Name sets the name in Contact.
func (c *Contact) Name(name string) *Contact {
	c.name = name
	return c
}

// Url sets the url in Contact.
func (c *Contact) Url(url string) *Contact {
	c.url = url
	return c
}

// Email sets the email in Contact.
func (c *Contact) Email(email string) *Contact {
	c.email = email
	return c
}

// Tag represents an api tag information of Document.
type Tag struct {
	name string
	desc string
}

// NewTag creates a default Tag with given arguments.
func NewTag(name, desc string) *Tag {
	return &Tag{name: name, desc: desc}
}

func (t *Tag) GetName() string { return t.name }
func (t *Tag) GetDesc() string { return t.desc }

// Name sets the name in Tag.
func (t *Tag) Name(name string) *Tag {
	t.name = name
	return t
}

// Desc sets the desc in Tag.
func (t *Tag) Desc(desc string) *Tag {
	t.desc = desc
	return t
}

// Security represents an api security information of Document.
type Security struct {
	title string
	typ   string // only supports apiKey
	in    string
	name  string
}

// NewSecurity creates a default Security with given arguments.
func NewSecurity(title, in, name string) *Security {
	return &Security{title: title, typ: "apiKey", in: in, name: name}
}

func (s *Security) GetTitle() string { return s.title }
func (s *Security) GetIn() string    { return s.in }
func (s *Security) GetName() string  { return s.name }

// Title sets the title in Security.
func (s *Security) Title(title string) *Security {
	s.title = title
	return s
}

// In sets the in in Security.
func (s *Security) In(in string) *Security {
	s.in = in
	return s
}

// Name sets the name in Security.
func (s *Security) Name(name string) *Security {
	s.name = name
	return s
}

// _document is the global Document.
var _document = NewDocument("", "", nil)

// SetDocument sets the basic information for global Document.
func SetDocument(host, basePath string, info *Info) *Document {
	return _document.Host(host).BasePath(basePath).Info(info)
}

// SetHost sets the host in global Document.
func SetHost(host string) *Document {
	return _document.Host(host)
}

// SetBasePath sets the basePath in global Document.
func SetBasePath(basePath string) *Document {
	return _document.BasePath(basePath)
}

// SetInfo sets the info in global Document.
func SetInfo(info *Info) *Document {
	return _document.Info(info)
}

// SetTags sets the whole tags in global Document.
func SetTags(tags ...*Tag) *Document {
	return _document.Tags(tags...)
}

// AddTags adds some tags into global Document.
func AddTags(tags ...*Tag) *Document {
	return _document.AddTags(tags...)
}

// SetSecurities sets the whole securities in global Document.
func SetSecurities(securities ...*Security) *Document {
	return _document.Securities(securities...)
}

// AddSecurities adds some securities into global Document.
func AddSecurities(securities ...*Security) *Document {
	return _document.AddSecurities(securities...)
}

// SetRoutePaths sets the whole route paths in global Document.
func SetRoutePaths(paths ...*RoutePath) *Document {
	return _document.RoutePaths(paths...)
}

// AddRoutePaths adds some route paths into global Document.
func AddRoutePaths(paths ...*RoutePath) *Document {
	return _document.AddRoutePaths(paths...)
}

// SetDefinitions sets the whole definitions in global Document.
func SetDefinitions(definitions ...*Definition) *Document {
	return _document.Definitions(definitions...)
}

// AddDefinitions adds some definitions into global Document.
func AddDefinitions(definitions ...*Definition) *Document {
	return _document.AddDefinitions(definitions...)
}
