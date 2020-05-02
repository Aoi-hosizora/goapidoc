package yamldoc

type Info struct {
	Title          string
	Description    string
	Version        string
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


type License struct {
	Name string
	Url  string
}

func NewLicense(name string, url string) *License {
	return &License{Name: name, Url: url}
}

type Contact struct {
	Name  string
	Url   string
	Email string
}

func NewContact(name string, url string, email string) *Contact {
	return &Contact{Name: name, Url: url, Email: email}
}

type Security struct {
	Title string
	Type  string // only support for apiKey
	Name  string
	In    string
}

func NewSecurity(title string, name string, in string) *Security {
	return &Security{Title: title, Type: "apiKey", Name: name, In: in}
}

type Tag struct {
	Name        string
	Description string
}

func NewTag(name string, description string) *Tag {
	return &Tag{Name: name, Description: description}
}
