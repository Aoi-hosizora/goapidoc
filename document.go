package apidoc

// Global document
var _document = &Document{}

func SetDocument(host string, basePath string, info *Info) {
	_document.Host = host
	_document.BasePath = basePath
	_document.Info = info
}

func SetTags(tags ...*Tag) {
	_document.SetTags(tags...)
}

func SetSecurities(securities ...*Security) {
	_document.SetSecurities(securities...)
}

func AddPaths(paths ...*Path) {
	_document.AddPaths(paths...)
}

func AddDefinitions(definitions ...*Definition) {
	_document.AddDefinitions(definitions...)
}

// Api document
type Document struct {
	Host     string
	BasePath string
	Info     *Info

	Tags        []*Tag
	Securities  []*Security
	Paths       []*Path
	Definitions []*Definition
}

func (d *Document) SetTags(tags ...*Tag) *Document {
	d.Tags = tags
	return d
}

func (d *Document) SetSecurities(security ...*Security) *Document {
	d.Securities = security
	return d
}

func (d *Document) AddPaths(path ...*Path) *Document {
	d.Paths = append(d.Paths, path...)
	return d
}

func (d *Document) AddDefinitions(models ...*Definition) *Document {
	d.Definitions = append(d.Definitions, models...)
	return d
}

// Base information
type Info struct {
	Title       string
	Description string
	Version     string

	TermsOfService string
	License        *License
	Contact        *Contact
}

func NewInfo(title string, description string, version string) *Info {
	return &Info{Title: title, Description: description, Version: version}
}

func (i *Info) SetTermsOfService(service string) *Info {
	i.TermsOfService = service
	return i
}

func (i *Info) SetLicense(license *License) *Info {
	i.License = license
	return i
}

func (i *Info) SetContact(contact *Contact) *Info {
	i.Contact = contact
	return i
}

// License information
type License struct {
	Name string
	Url  string
}

func NewLicense(name string, url string) *License {
	return &License{Name: name, Url: url}
}

// Contact information
type Contact struct {
	Name  string
	Url   string
	Email string
}

func NewContact(name string, url string, email string) *Contact {
	return &Contact{Name: name, Url: url, Email: email}
}

// Security information
type Security struct {
	Title string
	Type  string // only support for apiKey
	In    string
	Name  string
}

func NewSecurity(title string, in string, name string) *Security {
	return &Security{Title: title, Type: "apiKey", In: in, Name: name}
}

// Tag information
type Tag struct {
	Name        string
	Description string
}

func NewTag(name string, description string) *Tag {
	return &Tag{Name: name, Description: description}
}
