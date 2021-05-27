package goapidoc

// TODO BREAK CHANGES

// GenerateSwaggerYaml generates swagger yaml script and returns byte array.
func (d *Document) GenerateSwaggerYaml() ([]byte, error) {
	swagDoc := buildSwaggerDocument(d)
	return yamlMarshal(swagDoc)
}

// GenerateSwaggerJson generates swagger json script and returns byte array.
func (d *Document) GenerateSwaggerJson() ([]byte, error) {
	swagDoc := buildSwaggerDocument(d)
	return jsonMarshal(swagDoc)
}

// GenerateApib generates apib script and returns byte array.
func (d *Document) GenerateApib() ([]byte, error) {
	doc := buildApibDocument(d)
	return doc, nil
}

// SaveSwaggerYaml generates swagger yaml script and saves into file.
func (d *Document) SaveSwaggerYaml(path string) ([]byte, error) {
	bs, err := d.GenerateSwaggerYaml()
	if err != nil {
		return nil, err
	}
	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// SaveSwaggerJson generates swagger json script and saves into file.
func (d *Document) SaveSwaggerJson(path string) ([]byte, error) {
	bs, err := d.GenerateSwaggerJson()
	if err != nil {
		return nil, err
	}
	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// SaveApib generates apib script and saves into file.
func (d *Document) SaveApib(path string) ([]byte, error) {
	bs, err := d.GenerateApib()
	if err != nil {
		return nil, err
	}
	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// GenerateSwaggerYaml generates swagger yaml script and returns byte array.
func GenerateSwaggerYaml() ([]byte, error) {
	return _document.GenerateSwaggerYaml()
}

// GenerateSwaggerJson generates swagger json script and returns byte array.
func GenerateSwaggerJson() ([]byte, error) {
	return _document.GenerateSwaggerJson()
}

// GenerateApib generates apib script and returns byte array.
func GenerateApib() ([]byte, error) {
	return _document.GenerateApib()
}

// SaveSwaggerYaml generates swagger yaml script and saves into file.
func SaveSwaggerYaml(path string) ([]byte, error) {
	return _document.SaveSwaggerYaml(path)
}

// SaveSwaggerJson generates swagger json script and saves into file.
func SaveSwaggerJson(path string) ([]byte, error) {
	return _document.SaveSwaggerJson(path)
}

// SaveApib generates apib script and saves into file.
func SaveApib(path string) ([]byte, error) {
	return _document.SaveApib(path)
}
