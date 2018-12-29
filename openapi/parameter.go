package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Parameter Object
type Parameter struct {
	Name            string `yaml:"name,omitempty"`
	In              string `yaml:"in,omitempty"`
	Description     string `yaml:"description,omitempty"`
	Required        bool   `yaml:"required,omitempty"`
	Deprecated      string `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool   `yaml:"allowEmptyValue,omitempty"`

	Style         string              `yaml:"style,omitempty"`
	Explode       bool                `yaml:"explode,omitempty"`
	AllowReserved bool                `yaml:"allowReserved,omitempty"`
	Schema        *Schema             `yaml:"schema,omitempty"`
	Example       interface{}         `yaml:"example,omitempty"`
	Examples      map[string]*Example `yaml:"examples,omitempty"`

	Content map[string]*MediaType `yaml:"content,omitempty"`

	Ref string `yaml:"$ref,omitempty"`
}

// Validate the values of Parameter object.
// This function DOES NOT check whether the name field correspond to the associated path or not.
func (parameter Parameter) Validate() error {
	if parameter.Name == "" {
		return errors.New("parameter.name is required")
	}
	if parameter.In == "" {
		return errors.New("parameter.in is required")
	}
	if parameter.In == "path" && !parameter.Required {
		return errors.New("if parameter.in is path, required must be true")
	}
	validaters := []validater{parameter.Schema}
	if v, ok := parameter.Example.(validater); ok {
		validaters = append(validaters, v)
	}

	// example has no validation

	if len(parameter.Content) > 1 {
		return errors.New("parameter.content must only contain one entry")
	}
	for _, mediaType := range parameter.Content {
		validaters = append(validaters, mediaType)
	}
	return validateAll(validaters)
}
