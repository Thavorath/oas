package main

import (
	"github.com/thavorath/oas/openapi"
	"strings"
)

//AtRoute processes routes
func (v *Visitor) AtRoute(route *RouteParam) {
	fnc, src := route.FuncName, route.Src
	directives, _ := Parse(src)
	hasParams := false
	var pi *openapi.PathItem
	o := new(openapi.Operation)
	var method, path string
	for _, d := range directives {
		switch d.Name {
		case "route":
			method = d.Get("method")
			path = d.Get("path")
			pi = v.Doc.Paths[path]
			if pi == nil {
				pi = new(openapi.PathItem)
				v.Doc.Paths[path] = pi
			}
			o.Responses = map[string]*openapi.Response{}
			o.Description = d.Doc
			switch strings.ToUpper(method) {
			case "GET":
				pi.Get = o
			case "POST":
				pi.Post = o
			case "DELETE":
				pi.Delete = o
			case "PUT":
				pi.Put = o
			default:
				return
			}

			hasParams = (strings.Index(path, "{") != -1)
			for _, p := range d.Parameters {
				if p.Value == nil {
					o.Tags = append(o.Tags, p.Name)
				}
			}
		case "response":
			resp := new(openapi.Response)
			contentType := d.Get("content")
			if contentType == "" {
				contentType = "application/json"
			}
			resp.Content = map[string]*openapi.MediaType{}
			resp.Content[contentType] = &openapi.MediaType{
				Schema: &openapi.Schema{
					Ref: "#/components/schemas/" + d.Get("type"),
				},
			}
			resp.Description = v.GetSchemaDescription(d.Get("type"))
			o.Responses[d.Get("status")] = resp

		case "parameters":
			for _, p := range d.Parameters {
				prm := &openapi.Parameter{
					Ref:      "#/components/parameters/" + p.Name,
					
				}
				o.Parameters = append(o.Parameters, prm)
			}

		case "input":
			o.RequestBody = &openapi.RequestBody{
				Content: map[string]*openapi.MediaType{
					"application/json": &openapi.MediaType{
						Schema: &openapi.Schema{
							Ref: "#/components/schemas/" + d.Get("type"),
						},
					},
				},
				Required:    true,
				Description: v.GetSchemaDescription(d.Get("type")),
			}
		case "security":
			if d.Get("name") == "-" {
				o.Security = (*openapi.SecurityRequirement)(&[]map[string][]string{
					map[string][]string{},
				})
			} else {
				o.Security = (*openapi.SecurityRequirement)(&[]map[string][]string{
					map[string][]string{
						d.Get("name"): []string{},
					},
				})
			}

		}

	}

	o.OperationID = fnc
	o.Tags = append(o.Tags, route.Tag)
	if hasParams && len(o.Parameters) == 0 {
		println("Need parameters " + method + " " + path)
	}
}
