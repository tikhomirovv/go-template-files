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

## Usage

Place the `*.html` and `*.txt` template files in the same directory. File names of the same template (except extensions) must match. To get a template, you must specify the path to the directory and the name of the files relative to the selected directory with templates.

Consider an example. Let's say there is a templates directory at the root of the project that contains a welcome email template:

File `templates/greetings.html`:

```html
<html>
<head><title>{{title .Title}}</title></head>
<body><h1>Hello, {{.Username}}!</h1></body>
</html>
```

File `templates/greetings.txt`:

```txt
# {{title .Title}}

Hello, {{.Username}}!
```

Using the `//go:embed` directive, set the templates directory to search for templates. To get the `greetings` template, we will use the name `templates/greetings`.

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
	funcMap := templates.FuncMap{"title": strings.Title}
	vars := map[string]interface{}{"Title": "greetings!", "Username": "World"}

	// get template by path to template files
	tmpl := templates.Must(tmpls.Get("templates/greetings")).Funcs(funcMap)
    
	// apply a parsed template and write the output to wr
	var html, text bytes.Buffer
	if err := tmpl.Execute(&html, &text, vars); err != nil {
		panic(err)
	}

	fmt.Println(html.String())
	fmt.Println(text.String())
}
```

Output:

```html
<html>
<head><title>Greetings!</title></head>
<body><h1>Hello, World!</h1></body>
</html>
```

```txt
# Greetings!

Hello, World!
```


## TODO

 - [x] Parse template files and generate output 
 - [x] Use FuncMap for template engine
 - [ ] Automaticly parse and inject inline CSS
 - [ ] Minimize output
