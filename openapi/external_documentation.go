package openapi

// codebeat:disable[TOO_MANY_IVARS]

// ExternalDocumentation Object
type ExternalDocumentation struct {
	Description string `yaml:"description,omitempty"`
	URL         string `yaml:"url"`
}

// Validate the values of ExternalDocumentaion object.
func (externalDocumentation ExternalDocumentation) Validate() error {
	return mustURL("externalDocumentation.url", externalDocumentation.URL)
}
