package fylay

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateContext holds template variables and functions
type TemplateContext struct {
	variables map[string]interface{}
	funcs     template.FuncMap
}

// NewTemplateContext creates a new template context
func NewTemplateContext() *TemplateContext {
	return &TemplateContext{
		variables: make(map[string]interface{}),
		funcs:     make(template.FuncMap),
	}
}

// SetVariable sets a template variable
func (tc *TemplateContext) SetVariable(key string, value interface{}) {
	tc.variables[key] = value
}

// GetVariable gets a template variable
func (tc *TemplateContext) GetVariable(key string) (interface{}, bool) {
	val, ok := tc.variables[key]
	return val, ok
}

// SetFunc adds a custom template function
func (tc *TemplateContext) SetFunc(name string, fn interface{}) {
	tc.funcs[name] = fn
}

// ProcessTemplate processes a template string with the context
func (tc *TemplateContext) ProcessTemplate(tmplStr string) (string, error) {
	tmpl, err := template.New("layout").Funcs(tc.funcs).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tc.variables); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

// ProcessXMLWithTemplate processes XML content replacing template variables
func ProcessXMLWithTemplate(xmlContent string, ctx *TemplateContext) (string, error) {
	if ctx == nil {
		return xmlContent, nil
	}

	return ctx.ProcessTemplate(xmlContent)
}

// SetTemplateContext sets the template context on the builder
func (b *Builder) SetTemplateContext(ctx *TemplateContext) {
	b.templateContext = ctx
}

// GetTemplateContext returns the builder's template context
func (b *Builder) GetTemplateContext() *TemplateContext {
	if b.templateContext == nil {
		b.templateContext = NewTemplateContext()
	}
	return b.templateContext
}

// SetTemplateVariable sets a variable in the builder's template context
func (b *Builder) SetTemplateVariable(key string, value interface{}) {
	ctx := b.GetTemplateContext()
	ctx.SetVariable(key, value)
}

// ProcessLayoutTemplate processes the layout with template variables
func (b *Builder) ProcessLayoutTemplate(xmlContent []byte) ([]byte, error) {
	ctx := b.GetTemplateContext()
	if ctx == nil || len(ctx.variables) == 0 {
		return xmlContent, nil
	}

	processed, err := ProcessXMLWithTemplate(string(xmlContent), ctx)
	if err != nil {
		return nil, err
	}

	return []byte(processed), nil
}
