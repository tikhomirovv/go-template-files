# go-email-templates

[![GoDoc](https://godoc.org/github.com/tikhomirovv/go-email-templates?status.svg)](https://godoc.org/github.com/tikhomirovv/go-email-templates) [![Go Report Card](https://goreportcard.com/badge/github.com/tikhomirovv/go-email-templates)](https://goreportcard.com/report/github.com/tikhomirovv/go-email-templates)

<!-- [![GoCover](http://gocover.io/_badge/github.com/tikhomirovv/go-email-templates)](http://gocover.io/github.com/tikhomirovv/go-email-templates) -->

## Description

Simplifies the work with the conversion of email templates. Works with HTML and TXT files directly through the [is/fs](https://pkg.go.dev/io/fs) (to access files after compilation, it is recommended to use [embed](https://pkg.go.dev/embed) from standart library).

Powered by Go's [html/template](https://pkg.go.dev/html/template) and [text/template](https://pkg.go.dev/text/template) engine.

## Install

```sh
go get github.com/tikhomirovv/go-email-templates
```

## Basic usage

```go
import (
	templates "github.com/tikhomirovv/go-email-templates"
)

//go:embed templates
var templatesDir embed.FS

func main() {
	fsys := fs.FS(templatesDir)

	// Create default configuration
	cfg := templates.NewConfiguration(&fsys)
	tmpls := templates.NewTemplates(*cfg)

	// set funcMap & data variables
	funcMap := templates.FuncMap{"upper": strings.ToUpper}
	vars := map[string]interface{}{"Username": "Valerii"}

	// get template by path to template files
	tmpl := templates.Must(tmpls.Get("templates/email")).Funcs(funcMap)
    
	// apply a parsed template and write the output to wr
	var html, text bytes.Buffer
	if err := tmpl.Execute(&html, &text, vars); err != nil {
		panic(err)
	}

	fmt.Println(html.String())
	fmt.Println(text.String())
}
```

## TODO

 - [x] Parse template files and generate output 
 - [x] Use FuncMap for template engine
 - [ ] Automaticly parse and inject inline CSS
 - [ ] Minimize output
