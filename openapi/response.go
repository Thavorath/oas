package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Response Object
type Response struct {
	Description string                `yaml:"description,omitempty"`
	Headers     map[string]*Header    `yaml:"headers,omitempty"`
	Content     map[string]*MediaType `yaml:"content,omitempty"`
	Links       map[string]*Link      `yaml:"links,omitempty"`

	Ref string `yaml:"$ref,omitempty"`
}

// Validate the value of Response object.
func (response Response) Validate() error {
	if response.Description == "" {
		return errors.New("response.description is required")
	}
	validaters := []validater{}
	for _, header := range response.Headers {
		validaters = append(validaters, header)
	}
	for _, mediaType := range response.Content {
		validaters = append(validaters, mediaType)
	}
	for _, link := range response.Links {
		validaters = append(validaters, link)
	}
	return validateAll(validaters)
}
