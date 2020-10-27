package goapidoc

// type
const (
	INTEGER = "integer" // INTEGER type: integer, long
	NUMBER  = "number"  // NUMBER type: float, double
	STRING  = "string"  // STRING type: string, byte, binary, date, dateTime, password
	BOOLEAN = "boolean" // BOOLEAN type: boolean
	ARRAY   = "array"   // ARRAY type: array
	FILE    = "file"    // FILE type: file
	OBJECT  = "object"  // OBJECT x
)

// format
const (
	INT32    = "int32"     // INT32 format: signed 32 bits
	INT64    = "int64"     // INT64 format: signed 64 bits
	FLOAT    = "float"     // FLOAT format: float
	DOUBLE   = "double"    // DOUBLE format: double
	BYTE     = "byte"      // BYTE format: base64 encoded characters
	BINARY   = "binary"    // BINARY format: any sequence of octets
	DATE     = "date"      // DATE format: As defined by full-date - RFC3339
	DATETIME = "date-time" // DATETIME format: As defined by date-time - RFC3339
	PASSWORD = "password"  // PASSWORD format: Used to hint UIs the input needs to be obscured
)

// param
const (
	QUERY  = "query"    // QUERY param
	PATH   = "path"     // PATH param
	HEADER = "header"   // HEADER param
	BODY   = "body"     // BODY param
	FORM   = "formData" // FORM param
)

// method
const (
	GET     = "get"     // GET method
	PUT     = "put"     // PUT method
	POST    = "post"    // POST method
	DELETE  = "delete"  // DELETE method
	OPTIONS = "options" // OPTIONS method
	HEAD    = "head"    // HEAD method
	PATCH   = "patch"   // PATCH method
)

// mime
const (
	ALL   = "*/*"                               // ALL mime data: */*
	JSON  = "application/json"                  // JSON mime data: application/json
	XML   = "text/xml"                          // XML mime data: text/xml
	PLAIN = "text/plain"                        // PLAIN mime data: text/plain
	HTML  = "text/html"                         // HTML mime data: text/html
	MPFD  = "multipart/form-data"               // MPFD mime data: multipart/form-data
	URL   = "application/x-www-form-urlencoded" // URL mime data: application/x-www-form-urlencoded
	PNG   = "image/png"                         // PNG mime data: image/png
	JPEG  = "image/jpeg"                        // JPEG mime data: image/jpeg
	GIF   = "image/gif"                         // GIF mime data: image/gif
)
