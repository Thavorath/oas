package main

import (
	"fmt"
	"github.com/thavorath/oas/openapi"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"strings"
)

//Printf helper
func Printf(m string, v ...interface{}) {
	fmt.Printf(m, v...)
}

// Main entry point
func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	fs := token.NewFileSet()
	var v = new(Visitor)
	v.Doc.Version = "3.0.0"
	v.Doc.Info = new(openapi.Info)
	v.Doc.Paths = map[string]*openapi.PathItem{}
	v.Doc.ExternalDocs = &openapi.ExternalDocumentation{
		URL: "",
	}
	v.Doc.Components = &openapi.Components{
		Schemas:         map[string]*openapi.Schema{},
		Parameters:      map[string]*openapi.Parameter{},
		SecuritySchemes: map[string]*openapi.SecurityScheme{},
	}

	var parse func(string)
	parse = func(wd string) {
		pkgs, err := parser.ParseDir(fs, wd, nil, parser.ParseComments)
		for _, pkg := range pkgs {
			v.DefaultTag = pkg.Name
			ast.Walk(v, pkg)
		}
		f, err := os.Open(wd)
		if err != nil {
			return
		}
		defer f.Close()
		list, err := f.Readdir(-1)
		if err != nil {
			return
		}
		for _, fi := range list {
			if fi.IsDir() {
				parse(path.Join(wd, fi.Name()))
			}
		}
	}

	parse(cwd)
	for _, r := range v.Routes {
		v.AtRoute(r)
	}
	for _, m := range v.Merge {
		if s, ok := v.Doc.Components.Schemas[m.Name]; ok {
			for n, p := range s.Properties {
				m.Schema.Properties[n] = p
			}
		}
	}
	b, err := yaml.Marshal(v.Doc)
	if err != nil {
		fmt.Println("Error saving ", err)
		return
	}
	f, err := os.OpenFile("api.spec", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	f.Truncate(0)
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

type (
	//Visitor struct
	Visitor struct {
		DefaultTag string
		//holds the openapi document
		Doc    openapi.Document
		Routes []*RouteParam
		Merge  []MergeSpec
	}

	//MergeSpec struct
	MergeSpec struct {
		Schema *openapi.Schema
		Name   string
	}

	//RouteParam route param struct
	RouteParam struct {
		FuncName string
		Src      string
		Tag      string
	}
)

//GetSchemaDescription gets schema description
func (v *Visitor) GetSchemaDescription(name string) (description string) {
	if s, ok := v.Doc.Components.Schemas[name]; ok {
		description = s.Description
	}
	if description == "" {
		description = name //a default description
	}
	return
}

//Visit implements the visitor
func (v *Visitor) Visit(n ast.Node) ast.Visitor {
	if v == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.File:
		v.AtOpenapi(d.Doc.Text())
	case *ast.FuncDecl:
		v.Routes = append(v.Routes, &RouteParam{d.Name.Name, d.Doc.Text(), v.DefaultTag})
	case *ast.GenDecl:
		for _, s := range d.Specs {
			if st, ok := s.(*ast.TypeSpec); ok {

				cmt := d.Doc.Text()
				if cmt == "" {
					cmt = st.Doc.Text()
				}
				switch t := st.Type.(type) {
				case *ast.Ident:
					v.AtParameter(cmt, t)
				case *ast.StructType:
					if !strings.Contains(cmt, "@dto") {
						continue
					}
					v.AtDto(cmt, t)
				}
			}
		}
	}
	return v
}

func setType(schema *openapi.Schema, e ast.Expr, cmt string) {
	var array bool
	var t, t2 string
	params, customType := parseDirective(cmt, "@type:")
	if customType {
		t = params[0]
		if t == "map" {
			t2 = params[1]
		}
	} else {
	start:
		switch d2 := e.(type) {
		case *ast.Ident:
			t = d2.Name
		case *ast.StarExpr:
			e = d2.X
			goto start
		case *ast.ArrayType:
			array = true
			e = d2.Elt
			goto start
		case *ast.SelectorExpr:
			if fmt.Sprintf("%s.%s", d2.X, d2.Sel) == "time.Time" {
				t = "Time"
			}
		}
	}

	var s = schema
	if array {
		s.Type = "array"
		s = new(openapi.Schema)
		schema.Items = s
	}
	switch t {
	case "string":
		s.Type = "string"
	case "int", "int32", "int64":
		s.Type = "integer"
		s.Format = t
	case "float32":
		s.Type = "number"
		s.Format = "float"
	case "float64":
		s.Type = "number"
		s.Format = "double"
	case "bool":
		s.Type = "boolean"
	case "Time", "time":
		s.Type = "string"
		s.Format = "date-time"
	case "map":
		s.Type = "object"
		s.AdditionalProperties = &openapi.Schema{
			Type: t2,
		}
	case "base64":
		s.Type = "string"
		s.Format = "base64"
	default:
		s.Ref = "#/components/schemas/" + t
	}
}

func jsonTag(tag string) string {
	pos := strings.Index(tag, `json:"`)
	if pos == -1 {
		return ""
	}
	pos = pos + 6
	end := strings.Index(tag[pos:], `"`)
	if end == -1 {
		return ""
	}
	return strings.Split(tag[pos:pos+end], ",")[0]
}

func parseDirective(src, dir string) (params []string, found bool) {
	idx := strings.Index(src, dir)
	if idx == -1 {
		return nil, false
	}
	var line string
	end := strings.Index(src[idx:], "\n")
	if end == -1 {
		line = src[idx+len(dir):]
	} else {
		line = src[idx+len(dir) : end+idx]
	}
	return strings.Split(line, " "), true
}
