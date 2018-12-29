package openapi

import (
	"errors"
	"strconv"
)

// codebeat:disable[TOO_MANY_IVARS]

// Operation Object
type Operation struct {
	Tags         []string               `yaml:"tags,omitempty"`
	Summary      string                 `yaml:"summary,omitempty"`
	Description  string                 `yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty"`
	OperationID  string                 `yaml:"operationId,omitempty"`
	Parameters   []*Parameter           `yaml:"parameters,omitempty"`
	RequestBody  *RequestBody           `yaml:"requestBody,omitempty"`
	Responses    Responses              `yaml:"responses,omitempty"`
	Callbacks    map[string]*Callback   `yaml:"callbacks,omitempty"`
	Deprecated   bool                   `yaml:"deprecated,omitempty"`
	Security     *SecurityRequirement   `yaml:"security,omitempty"`
	Servers      []*Server              `yaml:"servers,omitempty"`
}

// SuccessResponse returns a success response object.
// If there are 2 or more success responses (like created and ok),
// it's not sure which is returned.
// If only match the default response or 2XX response, returned status code will be 0.
func (operation *Operation) SuccessResponse() (*Response, int, bool) {
	if operation == nil || operation.Responses == nil {
		return nil, -1, false
	}
	var defaultResponse *Response
	for statusStr, resp := range operation.Responses {
		switch statusStr {
		case "default":
			defaultResponse = resp
		case "2XX":
			defaultResponse = resp
		case "1XX", "3XX", "4XX", "5XX":
			continue
		}
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			continue
		}
		if statusInt/100 == 2 {
			if resp == nil {
				continue
			}
			return resp, statusInt, true
		}
	}
	return defaultResponse, 0, (defaultResponse != nil)
}

// Validate the values of Operation object.
func (operation Operation) Validate() error {
	validaters := []validater{}
	if operation.ExternalDocs != nil {
		validaters = append(validaters, operation.ExternalDocs)
	}
	if hasDuplicatedParameter(operation.Parameters) {
		return errors.New("some parameter is duplicated")
	}
	if operation.RequestBody != nil {
		validaters = append(validaters, operation.RequestBody)
	}
	if operation.Responses == nil {
		return errors.New("operation.responses is required")
	}
	validaters = append(validaters, operation.Responses)
	for _, callback := range operation.Callbacks {
		validaters = append(validaters, callback)
	}
	//for _, security := range operation.Security {
	//validaters = append(validaters, security)
	//}
	for _, server := range operation.Servers {
		validaters = append(validaters, server)
	}
	return validateAll(validaters)
}
