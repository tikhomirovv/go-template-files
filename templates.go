package templates

import (
	"bytes"

	htmlTemplate "html/template"
	textTemplate "text/template"
)

type Templates struct {
	cfg configuration
}

type TemplateName string

type Template struct {
	Name TemplateName
	Html string
	Text string
}

func NewTemplates(cfg configuration) *Templates {
	return &Templates{cfg}
}

func (t *Templates) GetTemplate(name TemplateName, vars interface{}) (*Template, error) {
	html, err := t.GetHtml(name, vars)
	if err != nil {
		return nil, err
	}

	text, err := t.GetText(name, vars)
	if err != nil {
		return nil, err
	}

	return &Template{
		Name: name,
		Html: html,
		Text: text,
	}, nil
}

func (t *Templates) GetHtml(name TemplateName, vars interface{}) (string, error) {
	fsys := *t.cfg.templatesFS
	templateFilename := string(name) + "." + t.cfg.Formats.Html.FileExtension
	templatesList, err := htmlTemplate.ParseFS(fsys, templateFilename)
	if err != nil {
		if t.cfg.Formats.Html.IsRequired {
			return "", FileNotFoundError(templateFilename)
		} else {
			return "", nil
		}
	}

	var content bytes.Buffer
	if err := templatesList.Execute(&content, vars); err != nil {
		return "", err
	}

	return content.String(), nil
}

func (t *Templates) GetText(name TemplateName, vars interface{}) (string, error) {
	fsys := *t.cfg.templatesFS
	templateFilename := string(name) + "." + t.cfg.Formats.Text.FileExtension
	templatesList, err := textTemplate.ParseFS(fsys, templateFilename)
	if err != nil {
		if t.cfg.Formats.Text.IsRequired {
			return "", FileNotFoundError(templateFilename)
		} else {
			return "", nil
		}
	}

	var content bytes.Buffer
	if err := templatesList.Execute(&content, vars); err != nil {
		return "", err
	}

	return content.String(), nil
}
