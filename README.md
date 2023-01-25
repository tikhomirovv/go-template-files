# go-template-files

[![GoDoc](https://godoc.org/github.com/tikhomirovv/go-template-files?status.svg)](https://godoc.org/github.com/tikhomirovv/go-template-files) [![Go Report Card](https://goreportcard.com/badge/github.com/tikhomirovv/go-template-files)](https://goreportcard.com/report/github.com/tikhomirovv/go-template-files)

<!-- [![GoCover](http://gocover.io/_badge/github.com/tikhomirovv/go-template-files)](http://gocover.io/github.com/tikhomirovv/go-template-files) -->

## Description

Simplifies work with template conversion. Just a wrapper for [html/template](https://pkg.go.dev/html/template) and [text/template](https://pkg.go.dev/text/template).
Works with HTML and TXT (any) files directly through the [is/fs](https://pkg.go.dev/io/fs) (to access files after compilation, it is recommended to use [embed](https://pkg.go.dev/embed) from standart library).

## Install

```sh
go get github.com/tikhomirovv/go-template-files
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
{{title .Title}}

Hello, {{.Username}}!
```

Using the `//go:embed` directive, set the templates directory to search for templates. To get the `greetings` template, we will use the name `templates/greetings`. By default, the configuration states that a file with `*.html` extenstion is required, but `*.txt` is not.

```go
import (
	ts "github.com/tikhomirovv/go-template-files"
)

//go:embed templates
var templatesDir embed.FS

func main() {
	fsys := fs.FS(templatesDir)

	// Create default configuration
	cfg := ts.NewConfiguration(&fsys)
	tmpls := ts.NewTemplates(*cfg)

	// set funcMap & data variables
	funcMap := ts.FuncMap{"title": strings.Title}
	vars := map[string]interface{}{"Title": "greetings!", "Username": "World"}

	// get template by path to template files
	tmpl := ts.Must(tmpls.Get("templates/greetings")).Funcs(funcMap)
    
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
Greetings!

Hello, World!
```

## Configuration

There are some configuration options available. For each of the html and txt formats, you can specify:

`FormatOptions.FileExtension` - what file extension to look for  
`FormatOptions.IsReguired` - whether an error will be thrown if a file with the specified extension is not found

For example, we want to process only the template markdown file:

File `templates/greetings.md`:

```md
# {{title .Title}}

Hello, *{{.Username}}*!
```

Set configuration:

```go
fsys := fs.FS(templatesDir)
cfg := ts.NewConfiguration(&fsys)
cfg.Formats[ts.Html].IsRequired = false
cfg.Formats[ts.Text].IsRequired = true
cfg.Formats[ts.Text].FileExtension = "md"

// or

cfg := &ts.Configuration{
	TemplatesFS: &fsys,
	Formats: ts.Formats{
		ts.Html: &ts.FormatOptions{
			IsRequired: false,
		},
		ts.Text: &ts.FormatOptions{
			FileExtension: "md",
			IsRequired:    true,
		},
	},
}
```

## TODO

 - [x] Parse template files and generate output 
 - [x] Use FuncMap for template engine
 - [ ] Automaticly parse and inject inline CSS
 - [ ] Minimize output
