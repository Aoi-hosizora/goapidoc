package goapidoc

// Document is the main api document.
type Document struct {
	host     string
	basePath string
	info     *Info

	tags        []*Tag
	securities  []*Security
	paths       []*RoutePath
	definitions []*Definition
}

// NewDocument creates a Document.
func NewDocument(host string, basePath string, info *Info) *Document {
	return &Document{host: host, basePath: basePath, info: info}
}

// Tags sets the tags in Document.
func (d *Document) Tags(tags ...*Tag) *Document {
	d.tags = tags
	return d
}

// Securities sets the securities in Document.
func (d *Document) Securities(security ...*Security) *Document {
	d.securities = security
	return d
}

// AddRoutePaths adds routePaths into Document.
func (d *Document) AddRoutePaths(path ...*RoutePath) *Document {
	d.paths = append(d.paths, path...)
	return d
}

// AddDefinitions adds definitions into Document.
func (d *Document) AddDefinitions(models ...*Definition) *Document {
	d.definitions = append(d.definitions, models...)
	return d
}

// Info is the basic information of Document.
type Info struct {
	title   string
	desc    string
	version string

	termsOfService string
	license        *License
	contact        *Contact
}

// NewInfo creates an Info.
func NewInfo(title string, desc string, version string) *Info {
	return &Info{title: title, desc: desc, version: version}
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

// License is the license information of Document.
type License struct {
	name string
	url  string
}

// NewLicense creates a License.
func NewLicense(name string, url string) *License {
	return &License{name: name, url: url}
}

// Contact is the contact information of Document.
type Contact struct {
	name  string
	url   string
	email string
}

// NewContact creates a Contact.
func NewContact(name string, url string, email string) *Contact {
	return &Contact{name: name, url: url, email: email}
}

// Security is the security information of Document.
type Security struct {
	title string
	typ   string // only support for apiKey
	in    string
	name  string
}

// NewSecurity creates a Security.
func NewSecurity(title string, in string, name string) *Security {
	return &Security{title: title, typ: "apiKey", in: in, name: name}
}

// Tag is the tag information of Document.
type Tag struct {
	name string
	desc string
}

// NewTag creates a Tag.
func NewTag(name string, desc string) *Tag {
	return &Tag{name: name, desc: desc}
}

// _document is the global Document.
var _document = NewDocument("", "", nil)

// SetDocument sets the global Document information.
func SetDocument(host string, basePath string, info *Info) {
	_document.host = host
	_document.basePath = basePath
	_document.info = info
}

// SetTags set the tags to global Document.
func SetTags(tags ...*Tag) {
	_document.Tags(tags...)
}

// SetSecurities set the securities to global Document.
func SetSecurities(securities ...*Security) {
	_document.Securities(securities...)
}

// AddRoutePaths adds routePaths into global Document.
func AddRoutePaths(paths ...*RoutePath) {
	_document.AddRoutePaths(paths...)
}

// AddDefinitions adds definitions into global Document.
func AddDefinitions(definitions ...*Definition) {
	_document.AddDefinitions(definitions...)
}
