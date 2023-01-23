# go-email-templates

## TODO

 - [x] Parse template files and generate output 
 - [x] Use FuncMap for template engine
 - [ ] Automaticly parse and inject inline CSS
 - [ ] Minimize output
 - [ ] Write output to files (?)
 - [ ] Preview file in browser command (?)


## Basic usage

```go
package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	templates "github.com/tikhomirovv/go-email-templates"
)

//go:embed templates
var templatesDir embed.FS

func main() {
	fsys := fs.FS(templatesDir)
	cfg := templates.NewConfiguration(&fsys)
	tmpls := templates.NewTemplates(*cfg)

	funcMap := templates.FuncMap{"upper": strings.ToUpper}
	vars := map[string]interface{}{"Username": "Valerii"}
	tmpl := templates.Must(tmpls.Get("templates/email")).Funcs(funcMap)
    
	var html bytes.Buffer
	if err := tmpl.ExecuteHtml(&html, vars); err != nil {
		panic(err)
	}

	fmt.Println(html.String())
}
```