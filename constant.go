package yamldoc

// noinspection GoUnusedConst
const (
	// type
	INTEGER = "integer" // type: integer, long
	NUMBER  = "number"  // type: float, double
	STRING  = "string"  // type: string, byte, binary, date, dateTime, password
	BOOLEAN = "boolean" // type: boolean
	OBJECT  = "object"  // type: customer
	ARRAY   = "array"   // type: customer

	// format
	INT32    = "int32"     // format: signed 32 bits
	INT64    = "int64"     // format: signed 64 bits
	FLOAT    = "float"     // format: float
	DOUBLE   = "double"    // format: double
	BYTE     = "byte"      // format: base64 encoded characters
	BINARY   = "binary"    // format: any sequence of octets
	DATE     = "date"      // format: As defined by full-date - RFC3339
	DATETIME = "date-time" // format: As defined by date-time - RFC3339
	PASSWORD = "password"  // format: Used to hint UIs the input needs to be obscured

	// param
	QUERY  = "query"    // param
	PATH   = "path"     // param
	HEADER = "header"   // param
	BODY   = "body"     // param
	FORM   = "formData" // param

	// mime
	JSON  = "application/json"                  // mime data
	XML   = "text/xml"                          // mime data
	PLAIN = "text/plain"                        // mime data
	HTML  = "text/html"                         // mime data
	MPFD  = "multipart/form-data"               // mime data
	URL   = "application/x-www-form-urlencoded" // mime data
	PNG   = "image/png"                         // mime data
	JPEG  = "image/jpeg"                        // mime data
	GIF   = "image/gif"                         // mime data
)
