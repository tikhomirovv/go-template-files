package templates

import (
	"io"
	"path"

	htmlTemplate "html/template"
	textTemplate "text/template"
)

type TemplateName string
type TemplateFilename string
type TemplateFiles map[Format]TemplateFilename
type FuncMap = textTemplate.FuncMap
type Templates struct {
	cfg Configuration
}
type Template struct {
	templates Templates
	Name      TemplateName
	files     TemplateFiles
	funcs     FuncMap
}

// Creates and configures a structure for working with templates
func NewTemplates(cfg Configuration) *Templates {
	return &Templates{cfg: cfg}
}

// SetConfiguration installs or updates a configuration
func (ts *Templates) SetConfiguration(cfg Configuration) {
	ts.cfg = cfg
}

// Get instantiates a template and looks for files that match the template.
func (ts *Templates) Get(path string) (*Template, error) {
	name := TemplateName(path)
	t := &Template{templates: *ts, Name: name}
	if err := t.getFiles(); err != nil {
		return nil, err
	}
	return t, nil
}

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable initializations
// such as
//
//	var t = templates.Must(tmplts.Get("templates/email")).Funcs(funcMap)
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

// getFiles scans the file system for files matching the template name and the formats
// described in the configuration. If the presence of a particular file format is marked
// as optional (`isRequired` == false) in the configuration, the file read error is suppressed
func (t *Template) getFiles() error {
	fsys := *t.templates.cfg.TemplatesFS
	formats := t.templates.cfg.Formats
	t.files = make(TemplateFiles)
	for format, opts := range formats {
		filename := string(t.Name) + "." + opts.FileExtension
		_, err := fsys.Open(filename)
		if err != nil {
			if !opts.IsRequired {
				continue
			}
			return err
		}
		t.files[format] = TemplateFilename(filename)
	}
	return nil
}

// Identical to the functions in packages `text/template` and `html/template`.
// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type. However, it is legal to overwrite elements of the map. The return
// value is the template, so calls can be chained.
func (t *Template) Funcs(funcMap FuncMap) *Template {
	t.funcs = funcMap
	return t
}

// Identical to the functions in package `html/template`.
// Execute applies a parsed HTML template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
func (t *Template) ExecuteHtml(wr io.Writer, vars interface{}) error {
	htmlFilename := t.files[Html]
	if htmlFilename == "" {
		return nil
	}
	template, err := htmlTemplate.
		New(path.Base(string(htmlFilename))).
		Funcs(t.funcs).
		ParseFS(*t.templates.cfg.TemplatesFS, string(htmlFilename))
	if err != nil {
		return err
	}
	return template.Funcs(t.funcs).Execute(wr, vars)
}

// Identical to the functions in package `text/template`.
// Execute applies a parsed text template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
func (t *Template) ExecuteText(wr io.Writer, vars interface{}) error {
	textFilename := t.files[Text]
	if textFilename == "" {
		return nil
	}
	template, err := textTemplate.
		New(path.Base(string(textFilename))).
		Funcs(t.funcs).
		ParseFS(*t.templates.cfg.TemplatesFS, string(textFilename))
	if err != nil {
		return err
	}
	return template.Execute(wr, vars)
}

// Execute applies a parsed HTML and text template to the specified data objects,
// writing the output to wr's.
func (t *Template) Execute(html io.Writer, text io.Writer, vars interface{}) error {
	if err := t.ExecuteHtml(html, vars); err != nil {
		return err
	}
	if err := t.ExecuteText(text, vars); err != nil {
		return err
	}
	return nil
}
