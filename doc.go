package yamldoc

type Document struct {
	Host     string
	BasePath string
	Info     *Info

	Tags     []*Tag
	Security []*Security
	Paths    []*Path
	Models   []*Model
}

func (d *Document) SetTags(tags ...*Tag) *Document {
	d.Tags = tags
	return d
}

func (d *Document) SetSecurities(security ...*Security) *Document {
	d.Security = security
	return d
}

func (d *Document) AddPath(path *Path) *Document {
	d.Paths = append(d.Paths, path)
	return d
}

func (d *Document) AddModel(model *Model) *Document {
	d.Models = append(d.Models, model)
	return d
}

var _document *Document

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

func AddPath(path *Path) {
	_document.AddPath(path)
}

func AddModel(model *Model) {
	_document.AddModel(model)
}
