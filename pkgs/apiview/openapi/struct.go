// package apiview

// // OpenApi Structs
// type OpenApiCollection struct {
// 	BasePath string                 `json:"basePath"`
// 	Host     string                 `json:"host"`
// 	Schemes  []string               `json:"schemes"`
// 	Protocol string                 `json:"protocol"`
// 	Info     OpenApiCollectionInfo  `json:"info"`
// 	Tags     []OpenApiTag           `json:"tags"`
// 	PathsRaw map[string]OpenApiPath `json:"paths"`
// 	Paths    []OpenApiPath          `json:"..."`
// }
//
// type OpenApiTag struct {
// 	Name        string
// 	Description string
// }
//
// type OpenApiCollectionInfo struct {
// 	Title       string `json:"title"`
// 	Version     string `json:"version"`
// 	Description string `json:"description"`
// 	Contact     struct {
// 		Email                   string `json:"email"`
// 		ResponsibleDeveloper    string `json:"responsibleDeveloper"`
// 		ResponsibleOrganization string `json:"responsibleOrganization"`
// 		Url                     string `json:"url"`
// 	}
// }
//
// type OpenApiPath struct {
// 	Method map[string]OpenApiEndpoint
// }
//
// type OpenApiEndpoint struct {
// 	Parameters []struct {
// 		In   string `json:"in"`
// 		Name string `json:"name"`
// 		Type string `json:"type"`
// 	} `json:"parameters"`
// 	Produces  []string                   `json:"produces"`
// 	Responses map[string]OpenApiResponse `json:"responses"`
// 	Summary   string                     `json:"summary"`
// 	Tags      []string                   `json:"tags"`
// }
//
// type OpenApiResponse struct {
// 	Description string `json:"description"`
// }

package openapi

// TODO: Check if fields are missing and used correct, missing 300 lines while unmarshaling and marshaling the httpbin.org OAS file

// OpenAPI represents the root of the OpenAPI document
type OpenAPI struct {
	OpenAPI      string                 `json:"openapi"`
	Info         Info                   `json:"info"`
	Paths        map[string]PathItem    `json:"paths"`
	Components   *Components            `json:"components,omitempty"`
	Servers      []Server               `json:"servers,omitempty"`
	Tags         []Tag                  `json:"tags,omitempty"`
	Security     []map[string][]string  `json:"security,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Info provides metadata about the API
type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Version        string   `json:"version"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Contact represents contact information
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License represents license information
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// PathItem represents an endpoint and its operations
type PathItem struct {
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty"`
	Post        *Operation  `json:"post,omitempty"`
	Put         *Operation  `json:"put,omitempty"`
	Delete      *Operation  `json:"delete,omitempty"`
	Options     *Operation  `json:"options,omitempty"`
	Head        *Operation  `json:"head,omitempty"`
	Patch       *Operation  `json:"patch,omitempty"`
	Trace       *Operation  `json:"trace,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// Operation represents a single API operation on a path
type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      interface{} `json:"schema,omitempty"` // JSON Schema object
}

// RequestBody represents the body of a request
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

// Response represents a response from an API operation
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Headers     map[string]Header    `json:"headers,omitempty"`
}

// MediaType represents a media type object
type MediaType struct {
	Schema   interface{}        `json:"schema,omitempty"` // JSON Schema object
	Example  interface{}        `json:"example,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
}

// Header represents a single header in a response
type Header struct {
	Description string      `json:"description,omitempty"`
	Schema      interface{} `json:"schema,omitempty"` // JSON Schema object
}

// Components contains reusable objects
type Components struct {
	Schemas       map[string]interface{} `json:"schemas,omitempty"` // JSON Schema objects
	Responses     map[string]Response    `json:"responses,omitempty"`
	Parameters    map[string]Parameter   `json:"parameters,omitempty"`
	RequestBodies map[string]RequestBody `json:"requestBodies,omitempty"`
}

// Server represents a server
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable represents a server variable for server URL template substitution
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// Tag represents a tag for API operations
type Tag struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// ExternalDocumentation represents external documentation
type ExternalDocumentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// Example represents an example object
type Example struct {
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty"`
}
