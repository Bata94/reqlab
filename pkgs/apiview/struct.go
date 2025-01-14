package apiview

type HTTPMethod int32

const (
	GET HTTPMethod = iota
	POST
	PUT
	DELETE
	PATCH
	HEAD
	OPTIONS
	CONNECT
	TRACE
)

var httpMethodName = map[HTTPMethod]string{
	GET:     "get",
	POST:    "post",
	PUT:     "put",
	DELETE:  "delete",
	PATCH:   "patch",
	HEAD:    "head",
	OPTIONS: "options",
	CONNECT: "connect",
	TRACE:   "trace",
}

func (h HTTPMethod) String() string {
	return httpMethodName[h]
}

type Endpoint struct {
	Path    string
	Method  HTTPMethod
	Headers []Header
}

type Header struct {
	Key     string
	Value   string
	Enabled bool
}
