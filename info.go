package main

import (
	"gitlab.com/thavorath/oas/openapi"
	"go/ast"
	"strings"
)

//AtOpenapi processes @openapi
func (v *Visitor) AtOpenapi(src string) {
	if !strings.HasPrefix(src, "@openapi") {
		return
	}

	directives, _ := Parse(src)
	for _, d := range directives {
		switch d.Name {
		case "openapi":
			v.Doc.Info.Title = d.Get("title")
			v.Doc.Info.Description = d.Doc
		case "version":
			v.Doc.Info.Version = d.Get("version")
		case "server":
			v.Doc.Servers = append(v.Doc.Servers, &openapi.Server{
				URL: d.Get("url"),
			})
		case "apiKey":
			v.Doc.Components.SecuritySchemes[d.Get("name")] = &openapi.SecurityScheme{
				Type: "apiKey",
				Name: d.Get("name"),
				In:   d.Get("in"),
			}
		case "security":
			v.Doc.Security = (*openapi.SecurityRequirement)(&[]map[string][]string{
				map[string][]string{
					d.Get("name"): []string{},
				},
			})
		}
	}

}

//AtParameter processes @parameter
func (v *Visitor) AtParameter(src string, t *ast.Ident) (ret bool) {
	if !strings.Contains(src, "@parameter:") {
		return false
	}
	rawLines := strings.Split(src, "\n")
	var lines []string
	var lineContainingParameter string
	for _, val := range rawLines {
		if strings.HasPrefix(val, "@parameter:") {
			lineContainingParameter = val
		} else {
			lines = append(lines, val)
		}
	}

	parts := strings.Split(lineContainingParameter, " ")
	if len(parts) != 3 {
		println("parameter should define a name and in. @parameter: customer_id path")
		return false
	}
	required := false
	if strings.HasPrefix(parts[1], "*") {
		required = true
		parts[1] = parts[1][1:]
	}
	p := &openapi.Parameter{
		Name:     parts[1],
		In:       parts[2],
		Required: required,
	}
	if p.Name == "" {
		return false
	}
	v.Doc.Components.Parameters[p.Name] = p
	for _, l := range lines {
		if strings.HasPrefix(l, "@required") {
			p.Required = true
		} else {
			p.Description = p.Description + "\n" + l
		}
	}
	if len(p.Description) > 0 {
		p.Description = p.Description[1:]
	}

	p.Schema = new(openapi.Schema)
	setType(p.Schema, t, p.Description)
	return true
}

//AtDto processes @dto
func (v *Visitor) AtDto(cmt string, t *ast.StructType) {
	schema := new(openapi.Schema)
	directives, _ := Parse(cmt)
	dtoName := ""
	for _, d := range directives {
		switch d.Name {
		case "dto":
			schema.Description = d.Doc
			v.Doc.Components.Schemas[d.Get("name")] = schema
			dtoName = d.Get("name")
		}
	}
	schema.Properties = map[string]*openapi.Schema{}
	for _, f := range t.Fields.List {
		if len(f.Names) == 0 {
			name := ""
			if sel, ok := f.Type.(*ast.SelectorExpr); ok {
				name = sel.Sel.Name
			} else if ident, ok := f.Type.(*ast.Ident); ok {
				name = ident.Name
			}
			v.Merge = append(v.Merge, MergeSpec{schema, name})
			continue
		}
		fs := new(openapi.Schema)
		name := f.Names[0].String()
		if !ast.IsExported(name) {
			continue
		}
		if f.Tag != nil {
			tag := jsonTag(f.Tag.Value)
			switch tag {
			case "":
			case "-":
				continue
			default:
				name = tag
			}
		} else {
			println("No tag " + f.Names[0].String() + " in type " + dtoName)
		}
		fs.Description = f.Doc.Text()
		if strings.HasPrefix(fs.Description, "*") {
			schema.Required = append(schema.Required, name)
			fs.Description = fs.Description[1:]
		}
		setType(fs, f.Type, fs.Description)
		schema.Properties[name] = fs
	}
}
