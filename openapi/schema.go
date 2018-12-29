package openapi

// codebeat:disable[TOO_MANY_IVARS]

// Schema Object
type Schema struct {
	Title            string   `yaml:"title,omitempty"`
	MultipleOf       int      `yaml:"multipleOf,omitempty"`
	Maximum          int      `yaml:"maximum,omitempty"`
	ExclusiveMaximum bool     `yaml:"exclusiveMaximum,omitempty"`
	Minimum          int      `yaml:"minimum,omitempty"`
	ExclusiveMinimum bool     `yaml:"exclusiveMinimum,omitempty"`
	MaxLength        int      `yaml:"maxLength,omitempty"`
	MinLength        int      `yaml:"minLength,omitempty"`
	Pattern          string   `yaml:"pattern,omitempty"`
	MaxItems         int      `yaml:"maxItems,omitempty"`
	MinItems         int      `yaml:"minItems,omitempty"`
	MaxProperties    int      `yaml:"maxProperties,omitempty"`
	MinProperties    int      `yaml:"minProperties,omitempty"`
	Required         []string `yaml:"required,omitempty"`
	Enum             []string `yaml:"enum,omitempty"`

	Type                 string             `yaml:"type,omitempty"`
	AllOf                *Schema            `yaml:"allOf,omitempty"`
	OneOf                *Schema            `yaml:"oneOf,omitempty"`
	AnyOf                *Schema            `yaml:"anyOf,omitempty"`
	Not                  *Schema            `yaml:"not,omitempty"`
	Items                *Schema            `yaml:"items,omitempty"`
	Properties           map[string]*Schema `yaml:"properties,omitempty"`
	AdditionalProperties *Schema            `yaml:"additionalProperties,omitempty"`
	Description          string             `yaml:"description,omitempty"`
	Format               string             `yaml:"format,omitempty"`
	Default              string             `yaml:"default,omitempty"`

	Nullable      bool                   `yaml:"nullable,omitempty"`
	Discriminator *Discriminator         `yaml:"descriminator,omitempty"`
	ReadOnly      bool                   `yaml:"readOnly,omitempty"`
	WriteOnly     bool                   `yaml:"writeOnly,omitempty"`
	XML           *XML                   `yaml:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `yaml:"externalDocs,omitempty"`
	Example       interface{}            `yaml:"example,omitempty"`
	Deprecated    bool                   `yaml:"deprecated,omitempty"`

	Ref string `yaml:"$ref,omitempty"`
}

// Validate the values of Schema object.
func (schema Schema) Validate() error {
	validaters := []validater{}
	if schema.AllOf != nil {
		validaters = append(validaters, schema.AllOf)
	}
	if schema.OneOf != nil {
		validaters = append(validaters, schema.OneOf)
	}
	if schema.AnyOf != nil {
		validaters = append(validaters, schema.AnyOf)
	}
	if schema.Not != nil {
		validaters = append(validaters, schema.Not)
	}
	if schema.Items != nil {
		validaters = append(validaters, schema.Items)
	}
	if schema.Discriminator != nil {
		validaters = append(validaters, schema.Discriminator)
	}
	if schema.XML != nil {
		validaters = append(validaters, schema.XML)
	}
	if schema.ExternalDocs != nil {
		validaters = append(validaters, schema.ExternalDocs)
	}
	for _, property := range schema.Properties {
		validaters = append(validaters, property)
	}
	if e, ok := schema.Example.(validater); ok {
		validaters = append(validaters, e)
	}
	return validateAll(validaters)
}
