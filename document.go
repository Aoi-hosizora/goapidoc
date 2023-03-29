package goapidoc

// ========
// Document
// ========

// Document represents an api document information.
type Document struct {
	host     string
	basePath string
	info     *Info

	option      *Option
	operations  []*Operation
	definitions []*Definition
}

// NewDocument creates a default Document with given arguments.
func NewDocument(host, basePath string, info *Info) *Document {
	return &Document{host: host, basePath: basePath, info: info}
}

// GetHost returns the host from Document.
func (d *Document) GetHost() string { return d.host }

// GetBasePath returns the basePath from Document.
func (d *Document) GetBasePath() string { return d.basePath }

// GetInfo returns the info from Document.
func (d *Document) GetInfo() *Info { return d.info }

// GetOption returns the option from Document.
func (d *Document) GetOption() *Option { return d.option }

// GetOperations returns the whole operations from Document.
func (d *Document) GetOperations() []*Operation { return d.operations }

// GetDefinitions returns the whole definitions from Document.
func (d *Document) GetDefinitions() []*Definition { return d.definitions }

// Cleanup cleans up the Document.
func (d *Document) Cleanup() *Document {
	d.host = ""
	d.basePath = ""
	d.info = nil
	d.option = nil
	d.operations = nil
	d.definitions = nil
	return d
}

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

// Option sets the option in Document.
func (d *Document) Option(option *Option) *Document {
	d.option = option
	return d
}

// Operations sets the whole operations in Document.
func (d *Document) Operations(operations ...*Operation) *Document {
	d.operations = operations
	return d
}

// AddOperations adds some operations into Document.
func (d *Document) AddOperations(operations ...*Operation) *Document {
	d.operations = append(d.operations, operations...)
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

// ====
// Info
// ====

// Info represents a basic information of Document.
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

// GetTitle returns the title from Info.
func (i *Info) GetTitle() string { return i.title }

// GetDesc returns the desc from Info.
func (i *Info) GetDesc() string { return i.desc }

// GetVersion returns the version from Info.
func (i *Info) GetVersion() string { return i.version }

// GetTermsOfService returns the terms of service from Info.
func (i *Info) GetTermsOfService() string { return i.termsOfService }

// GetLicense returns the license from Info.
func (i *Info) GetLicense() *License { return i.license }

// GetContact returns the contact from Info.
func (i *Info) GetContact() *Contact { return i.contact }

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

// TermsOfService sets the terms of service in Info.
func (i *Info) TermsOfService(terms string) *Info {
	i.termsOfService = terms
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

// =======
// License
// =======

// License represents a license information of Document.
type License struct {
	name string
	url  string
}

// NewLicense creates a default License with given arguments.
func NewLicense(name, url string) *License {
	return &License{name: name, url: url}
}

// GetName returns the name from License.
func (l *License) GetName() string { return l.name }

// GetUrl returns the url from License.
func (l *License) GetUrl() string { return l.url }

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

// =======
// Contact
// =======

// Contact represents a contact information of Document.
type Contact struct {
	name  string
	url   string
	email string
}

// NewContact creates a default Contact with given arguments.
func NewContact(name, url, email string) *Contact {
	return &Contact{name: name, url: url, email: email}
}

// GetName returns the name from Contact.
func (c *Contact) GetName() string { return c.name }

// GetUrl returns the url from Contact.
func (c *Contact) GetUrl() string { return c.url }

// GetEmail returns the email from Contact.
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

// ======
// Option
// ======

// Option represents an api extra option of Document.
type Option struct {
	schemes       []string
	consumes      []string
	produces      []string
	tags          []*Tag
	securities    []*Security
	globalParams  []*Param
	externalDoc   *ExternalDoc
	additionalDoc string
	routesOptions []*RoutesOption
}

// NewOption creates a default Option.
func NewOption() *Option {
	return &Option{}
}

// GetSchemes returns the whole schemes from Option.
func (o *Option) GetSchemes() []string { return o.schemes }

// GetConsumes returns the whole consumes from Option.
func (o *Option) GetConsumes() []string { return o.consumes }

// GetProduces returns the whole produces from Option.
func (o *Option) GetProduces() []string { return o.produces }

// GetTags returns the whole tags from Option.
func (o *Option) GetTags() []*Tag { return o.tags }

// GetSecurities returns the whole securities from Option.
func (o *Option) GetSecurities() []*Security { return o.securities }

// GetGlobalParams returns the whole global params from Option.
func (o *Option) GetGlobalParams() []*Param { return o.globalParams }

// GetExternalDoc returns the external documentation from Option.
func (o *Option) GetExternalDoc() *ExternalDoc { return o.externalDoc }

// GetAdditionalDoc returns the additional document from Option.
func (o *Option) GetAdditionalDoc() string { return o.additionalDoc }

// GetRoutesOptions returns the whole routes options from Option.
func (o *Option) GetRoutesOptions() []*RoutesOption { return o.routesOptions }

// Schemes sets the whole schemes in Option.
func (o *Option) Schemes(schemes ...string) *Option {
	o.schemes = schemes
	return o
}

// AddSchemes adds some schemes into Option.
func (o *Option) AddSchemes(schemes ...string) *Option {
	o.schemes = append(o.schemes, schemes...)
	return o
}

// Consumes sets the whole consumes in Option.
func (o *Option) Consumes(consumes ...string) *Option {
	o.consumes = consumes
	return o
}

// AddConsumes adds some consumes into Option.
func (o *Option) AddConsumes(consumes ...string) *Option {
	o.consumes = append(o.consumes, consumes...)
	return o
}

// Produces sets the whole produces in Option.
func (o *Option) Produces(produces ...string) *Option {
	o.produces = produces
	return o
}

// AddProduces adds some produces into Option.
func (o *Option) AddProduces(produces ...string) *Option {
	o.produces = append(o.produces, produces...)
	return o
}

// Tags sets the whole tags in Option.
func (o *Option) Tags(tags ...*Tag) *Option {
	o.tags = tags
	return o
}

// AddTags adds some tags into Option.
func (o *Option) AddTags(tags ...*Tag) *Option {
	o.tags = append(o.tags, tags...)
	return o
}

// Securities sets the whole securities in Option.
func (o *Option) Securities(securities ...*Security) *Option {
	o.securities = securities
	return o
}

// AddSecurities adds some securities into Option.
func (o *Option) AddSecurities(securities ...*Security) *Option {
	o.securities = append(o.securities, securities...)
	return o
}

// GlobalParams sets the whole global params in Option.
func (o *Option) GlobalParams(globalParams ...*Param) *Option {
	o.globalParams = globalParams
	return o
}

// AddGlobalParams adds some global params into Option.
func (o *Option) AddGlobalParams(globalParams ...*Param) *Option {
	o.globalParams = append(o.globalParams, globalParams...)
	return o
}

// ExternalDoc sets the external documentation in Option.
func (o *Option) ExternalDoc(doc *ExternalDoc) *Option {
	o.externalDoc = doc
	return o
}

// AdditionalDoc sets the additional document in Option, this is only supported in API Blueprint.
func (o *Option) AdditionalDoc(doc string) *Option {
	o.additionalDoc = doc
	return o
}

// RoutesOptions sets the whole routes group options in Option, this is only supported in API Blueprint.
func (o *Option) RoutesOptions(options ...*RoutesOption) *Option {
	o.routesOptions = options
	return o
}

// AddRoutesOptions adds some routes group options in Option, this is only supported in API Blueprint.
func (o *Option) AddRoutesOptions(options ...*RoutesOption) *Option {
	o.routesOptions = append(o.routesOptions, options...)
	return o
}

// ===
// Tag
// ===

// Tag represents a tag information of Document.
type Tag struct {
	name          string
	desc          string
	externalDoc   *ExternalDoc
	additionalDoc string
}

// NewTag creates a default Tag with given arguments.
func NewTag(name, desc string) *Tag {
	return &Tag{name: name, desc: desc}
}

// GetName returns the name from Tag.
func (t *Tag) GetName() string { return t.name }

// GetDesc returns the desc from Tag.
func (t *Tag) GetDesc() string { return t.desc }

// GetExternalDoc returns the external documentation from Tag.
func (t *Tag) GetExternalDoc() *ExternalDoc { return t.externalDoc }

// GetAdditionalDoc returns the additional document from Tag.
func (t *Tag) GetAdditionalDoc() string { return t.additionalDoc }

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

// ExternalDoc sets the external documentation in Tag.
func (t *Tag) ExternalDoc(doc *ExternalDoc) *Tag {
	t.externalDoc = doc
	return t
}

// AdditionalDoc sets the additional document in Tag, this is only supported in API Blueprint.
func (t *Tag) AdditionalDoc(doc string) *Tag {
	t.additionalDoc = doc
	return t
}

// ========
// Security
// ========

// Security represents a security definition information of Document.
type Security struct {
	title string
	typ   string
	desc  string

	in   string // only for apiKey
	name string // only for apiKey

	flow             string           // only for oauth2
	authorizationUrl string           // only for oauth2
	tokenUrl         string           // only for oauth2
	scopes           []*SecurityScope // only for oauth2
}

// NewSecurity creates a default Security with given arguments.
func NewSecurity(title string, typ string) *Security {
	return &Security{title: title, typ: typ}
}

// NewApiKeySecurity creates an apiKey authentication Security with given arguments.
func NewApiKeySecurity(title, in, name string) *Security {
	return &Security{title: title, typ: APIKEY, in: in, name: name}
}

// NewBasicSecurity creates a basic authentication Security with given arguments.
func NewBasicSecurity(title string) *Security {
	return &Security{title: title, typ: BASIC}
}

// NewOAuth2Security creates an oauth2 authentication Security with given arguments.
func NewOAuth2Security(title string, flow string) *Security {
	return &Security{title: title, typ: OAUTH2, flow: flow}
}

// GetTitle returns the title from Security.
func (s *Security) GetTitle() string { return s.title }

// GetType returns the authentication type from Security.
func (s *Security) GetType() string { return s.typ }

// GetDesc returns the desc from Security.
func (s *Security) GetDesc() string { return s.desc }

// GetInLoc returns the in-location from Security.
func (s *Security) GetInLoc() string { return s.in }

// GetName returns the name from Security.
func (s *Security) GetName() string { return s.name }

// GetFlow returns the flow from Security.
func (s *Security) GetFlow() string { return s.flow }

// GetAuthorizationUrl returns the authorization url from Security.
func (s *Security) GetAuthorizationUrl() string { return s.authorizationUrl }

// GetTokenUrl returns the token url from Security.
func (s *Security) GetTokenUrl() string { return s.tokenUrl }

// GetScopes returns the whole security scopes from Security.
func (s *Security) GetScopes() []*SecurityScope { return s.scopes }

// Title sets the title in Security.
func (s *Security) Title(title string) *Security {
	s.title = title
	return s
}

// Type sets the authentication type in Security.
func (s *Security) Type(typ string) *Security {
	s.typ = typ
	return s
}

// Desc sets the desc in Security.
func (s *Security) Desc(desc string) *Security {
	s.desc = desc
	return s
}

// InLoc sets the in-location in Security.
func (s *Security) InLoc(in string) *Security {
	s.in = in
	return s
}

// Name sets the name in Security.
func (s *Security) Name(name string) *Security {
	s.name = name
	return s
}

// Flow sets the flow in Security.
func (s *Security) Flow(flow string) *Security {
	s.flow = flow
	return s
}

// AuthorizationUrl sets the authorization url in Security.
func (s *Security) AuthorizationUrl(authorizationUrl string) *Security {
	s.authorizationUrl = authorizationUrl
	return s
}

// TokenUrl sets the token url in Security.
func (s *Security) TokenUrl(tokenUrl string) *Security {
	s.tokenUrl = tokenUrl
	return s
}

// Scopes sets the whole security scopes in Security.
func (s *Security) Scopes(scopes ...*SecurityScope) *Security {
	s.scopes = scopes
	return s
}

// AddScopes add some security scopes into Security.
func (s *Security) AddScopes(scopes ...*SecurityScope) *Security {
	s.scopes = append(s.scopes, scopes...)
	return s
}

// =============
// SecurityScope
// =============

// SecurityScope represents a security scope information of Document.
type SecurityScope struct {
	scope string
	desc  string
}

// NewSecurityScope creates a default SecurityScope with given arguments.
func NewSecurityScope(scope, desc string) *SecurityScope {
	return &SecurityScope{scope: scope, desc: desc}
}

// GetScope returns the scope from SecurityScope.
func (s *SecurityScope) GetScope() string { return s.scope }

// GetDesc returns the desc from SecurityScope.
func (s *SecurityScope) GetDesc() string { return s.desc }

// Scope sets the scope in Security.
func (s *SecurityScope) Scope(scope string) *SecurityScope {
	s.scope = scope
	return s
}

// Desc sets the desc in Security.
func (s *SecurityScope) Desc(desc string) *SecurityScope {
	s.desc = desc
	return s
}

// ===========
// ExternalDoc
// ===========

// ExternalDoc represents an external documentation information of Document, Tag and Operation.
type ExternalDoc struct {
	desc string
	url  string
}

// NewExternalDoc creates a default ExternalDoc with given arguments.
func NewExternalDoc(desc, url string) *ExternalDoc {
	return &ExternalDoc{desc: desc, url: url}
}

// GetDesc returns the desc from ExternalDoc.
func (e *ExternalDoc) GetDesc() string { return e.desc }

// GetUrl returns the url from ExternalDoc.
func (e *ExternalDoc) GetUrl() string { return e.url }

// Desc sets the desc in ExternalDoc.
func (e *ExternalDoc) Desc(desc string) *ExternalDoc {
	e.desc = desc
	return e
}

// Url sets the url in ExternalDoc.
func (e *ExternalDoc) Url(url string) *ExternalDoc {
	e.url = url
	return e
}

// ============
// RoutesOption
// ============

// RoutesOption represents a routes group option of Document, this is only supported in API Blueprint.
type RoutesOption struct {
	route         string
	summary       string
	additionalDoc string
}

// NewRoutesOption creates a default RoutesOption with given arguments.
func NewRoutesOption(route string) *RoutesOption {
	return &RoutesOption{route: route}
}

// GetRoute returns the route from RoutesOption.
func (r *RoutesOption) GetRoute() string { return r.route }

// GetSummary returns the summary from RoutesOption.
func (r *RoutesOption) GetSummary() string { return r.summary }

// GetAdditionalDoc returns the additional document from RoutesOption.
func (r *RoutesOption) GetAdditionalDoc() string { return r.additionalDoc }

// Route sets the route in RoutesOption.
func (r *RoutesOption) Route(route string) *RoutesOption {
	r.route = route
	return r
}

// Summary sets the summary in RoutesOption.
func (r *RoutesOption) Summary(summary string) *RoutesOption {
	r.summary = summary
	return r
}

// AdditionalDoc sets the additional document in RoutesOption, this is only supported in API Blueprint.
func (r *RoutesOption) AdditionalDoc(additionalDoc string) *RoutesOption {
	r.additionalDoc = additionalDoc
	return r
}

// ===============
// global document
// ===============

// _document represents a global Document.
var _document = NewDocument("", "", nil)

// SetDocument sets the basic information for global Document.
func SetDocument(host, basePath string, info *Info) *Document {
	_document.host = host
	_document.basePath = basePath
	_document.info = info
	return _document
}

// GetHost returns the host from global Document.
func GetHost() string { return _document.GetHost() }

// GetBasePath returns the basePath from global Document.
func GetBasePath() string { return _document.GetBasePath() }

// GetInfo returns the info from global Document.
func GetInfo() *Info { return _document.GetInfo() }

// GetOption returns the option from global Document.
func GetOption() *Option { return _document.GetOption() }

// GetOperations returns the whole operations from global Document.
func GetOperations() []*Operation { return _document.GetOperations() }

// GetDefinitions returns the whole definitions from global Document.
func GetDefinitions() []*Definition { return _document.GetDefinitions() }

// CleanupDocument cleans up the global Document.
func CleanupDocument() *Document {
	return _document.Cleanup()
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

// SetOption sets the option in Document.
func SetOption(option *Option) *Document {
	return _document.Option(option)
}

// SetOperations sets the whole operations in global Document.
func SetOperations(operations ...*Operation) *Document {
	return _document.Operations(operations...)
}

// AddOperations adds some operations into global Document.
func AddOperations(operations ...*Operation) *Document {
	return _document.AddOperations(operations...)
}

// SetDefinitions sets the whole definitions in global Document.
func SetDefinitions(definitions ...*Definition) *Document {
	return _document.Definitions(definitions...)
}

// AddDefinitions adds some definitions into global Document.
func AddDefinitions(definitions ...*Definition) *Document {
	return _document.AddDefinitions(definitions...)
}
