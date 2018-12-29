package main

import (
	"strings"
	"text/scanner"
)

var constDirectives = map[string][]string{
	"openapi":  []string{"title"},
	"version":  []string{"version"},
	"server":   []string{"url"},
	"apiKey":   []string{"in", "name"},
	"security": []string{"name"},

	"dto":       []string{"name"},
	"parameter": []string{"name"},
	"type":      []string{"type"},

	"route":      []string{"method", "path"},
	"parameters": []string{},
	"response":   []string{"status", "type"},
	"input":      []string{"type"},
}

//Directive represents a directive
type Directive struct {
	Name       string
	Parameters []Parameter
	Doc        string
}

//Parameter represents a parameter
type Parameter struct {
	Name  string
	Value *string
}

//Get gets param value
func (d *Directive) Get(name string) string {
	for _, p := range d.Parameters {
		if p.Name == name {
			return *p.Value
		}
	}
	return ""
}

//Parse parses string and extracts directives
func Parse(src string) (directives []*Directive, err error) {
	lines := strings.Split(src, "\n")
	var d *Directive
	var args []string
	var initialDocs string
	for _, l := range lines {
		l = strings.Trim(l, " \t")
		if !strings.HasPrefix(l, "@") {
			if d != nil {
				d.Doc = d.Doc + "\n" + l
			} else {
				initialDocs = initialDocs + "\n" + l
			}
			continue
		}

		idx := strings.Index(l, ":")
		if idx == -1 {
			continue
		}
		name := l[1:idx]
		if _, ok := constDirectives[name]; !ok {
			continue
		}
		d = new(Directive)
		d.Name = name
		if initialDocs != "" {
			d.Doc = initialDocs
			initialDocs = ""
		}
		args = constDirectives[d.Name]
		var s scanner.Scanner
		s.Init(strings.NewReader(l[idx+1:]))
		foundEqual := false
		for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
			txt := strings.Trim(s.TokenText(), `'"`)
			if foundEqual && len(d.Parameters) > 0 {
				d.Parameters[len(d.Parameters)-1].Value = &txt
				foundEqual = false
				continue
			}
			if txt == "=" {
				foundEqual = true
			} else {
				d.Parameters = append(d.Parameters, Parameter{
					Name: txt,
				})
				if len(d.Parameters) > 1 && d.Parameters[len(d.Parameters)-2].Value == nil && len(args) > 0 {
					p := &d.Parameters[len(d.Parameters)-2]
					name := p.Name
					p.Value = &name
					p.Name = args[0]
					args = args[1:]
				}
			}
		}

		if len(d.Parameters) > 0 && d.Parameters[len(d.Parameters)-1].Value == nil && len(args) > 0 {
			p := &d.Parameters[len(d.Parameters)-1]
			name := p.Name
			p.Value = &name
			p.Name = args[0]
		}

		directives = append(directives, d)
	}

	for _, d := range directives {
		d.Doc = strings.Trim(d.Doc, " \n")
	}
	return directives, err
}
