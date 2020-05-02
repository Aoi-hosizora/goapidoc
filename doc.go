package yamldoc

type Document struct {
	Host     string
	BasePath string
	Info     *Info

	Tags       []*Tag
	Securities []*Security
	Paths      []*Path
	Models     []*Model
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

func (d *Document) AddModels(models ...*Model) *Document {
	d.Models = append(d.Models, models...)
	return d
}

// global document
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

func AddModels(models ...*Model) {
	_document.AddModels(models...)
}
