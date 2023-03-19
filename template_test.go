package templates_test

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ts "github.com/tikhomirovv/go-template-files"
)

type file struct {
	name     string
	contents string
}

// createFS creates template files in temporary directory
func createFS(t *testing.T, files []file) fs.FS {
	td := t.TempDir()
	for _, file := range files {
		err := os.MkdirAll(filepath.Join(td, filepath.Dir(file.name)), 0750)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(td, file.name), []byte(file.contents), 0600)
		if err != nil {
			t.Fatal(err)
		}
	}

	return os.DirFS(td)
}

// TestGetTemplates tests getting access to template files depending on the configuration
func TestGetTemplates(t *testing.T) {
	fsys := createFS(t, []file{
		// check parse templates
		{"template1/one.html", `HTML-code`},
		{"template1/one.txt", `Plain text`},

		// the template have html only, no txt
		{"template2/sub/two.html", `-`},
	})
	cfg := ts.NewConfiguration(&fsys)
	tpls := ts.NewTemplates(*cfg)

	t.Run("HTML:required", func(t *testing.T) {
		cfg.Formats[ts.Html].IsRequired = true
		tpls.SetConfiguration(*cfg)

		tmpl, err := tpls.Get("template1/no-template")
		testTemplateNotExists(t, tmpl, err)

		tmpl, err = tpls.Get("template1/one")
		testTemplateExists(t, tmpl, err)

		tmpl, err = tpls.Get("template2/sub/two")
		testTemplateExists(t, tmpl, err)
	})

	t.Run("HTML:optional", func(t *testing.T) {
		cfg.Formats[ts.Html].IsRequired = false
		tpls.SetConfiguration(*cfg)
		tmpl, err := tpls.Get("template1/no-template")
		testTemplateExists(t, tmpl, err)
	})

	t.Run("TXT:required", func(t *testing.T) {
		cfg.Formats[ts.Text].IsRequired = true
		tpls.SetConfiguration(*cfg)

		tmpl, err := tpls.Get("template2/subdir/two")
		testTemplateNotExists(t, tmpl, err)

		tmpl, err = tpls.Get("template1/one")
		testTemplateExists(t, tmpl, err)
	})

	t.Run("TXT:optional", func(t *testing.T) {
		cfg.Formats[ts.Text].IsRequired = false
		tpls.SetConfiguration(*cfg)
		tmpl, err := tpls.Get("template2/subdir/two")
		testTemplateExists(t, tmpl, err)
	})
}

// testTemplateExists makes sure the template is accessed
func testTemplateExists(t *testing.T, tmpl *ts.Template, err error) {
	if err != nil {
		t.Error(err)
	}

	if tmpl == nil {
		t.Error("the template must be found")
	}
}

// testTemplateNotExists makes sure the template is not accessed
func testTemplateNotExists(t *testing.T, tmpl *ts.Template, err error) {
	if tmpl != nil {
		t.Error("the template must not be found")
	}
}

// TestExecute tests template conversion to final result using functions and variables
func TestExecute(t *testing.T) {
	templateName := "template2/sub/func"
	fsys := createFS(t, []file{
		// check funcMap
		{templateName + ".html", `{{template "header.html" .}}<h1>Hello, {{.Username}}</h1>{{template "footer.html" .}}`},
		{templateName + ".txt", `# {{upper .Title}}. {{template "additional.txt" .}}`},

		// common templates
		{"common/header.html", `<html><head><title>{{upper .Title}}</title></head><body>`},
		{"common/footer.html", `</body></html>`},
		{"common/additional.txt", `Help me please.`},
	})
	cfg := ts.NewConfiguration(&fsys)
	cfg.SetCommonTemplatesPath("common")
	tmpls := ts.NewTemplates(*cfg)
	funcMap := ts.FuncMap{"upper": strings.ToUpper}
	vars := map[string]interface{}{"Title": "Hello from tests", "Username": "Valerii"}
	tmpl := ts.Must(tmpls.Get(templateName)).Funcs(funcMap)

	expectHtml := `<html><head><title>HELLO FROM TESTS</title></head><body><h1>Hello, Valerii</h1></body></html>`
	expectText := `# HELLO FROM TESTS. Help me please.`

	var html, text bytes.Buffer

	t.Run("HTML:execute", func(t *testing.T) {
		if err := tmpl.ExecuteHtml(&html, vars); err != nil {
			t.Fatal(err)
		}
		gotHtml := html.String()
		if gotHtml != expectHtml {
			t.Errorf("incorrect html template processing; expect %#q; got: %#q", expectHtml, gotHtml)
		}
	})

	t.Run("TXT:execute", func(t *testing.T) {
		if err := tmpl.ExecuteText(&text, vars); err != nil {
			t.Fatal(err)
		}
		gotText := text.String()
		if gotText != expectText {
			t.Errorf("incorrect text template processing; expect %#q; got: %#q", expectText, gotText)
		}
	})

	t.Run("ALL:execute", func(t *testing.T) {
		html.Reset()
		text.Reset()
		if err := tmpl.Execute(&html, &text, vars); err != nil {
			t.Fatal(err)
		}
		gotHtml := html.String()
		gotText := text.String()
		if gotHtml != expectHtml {
			t.Errorf("incorrect html template processing; expect %#q; got: %#q", expectHtml, gotHtml)
		}
		if gotText != expectText {
			t.Errorf("incorrect text template processing; expect %#q; got: %#q", expectText, gotText)
		}
	})
}

func TestSetConfiguration(t *testing.T) {
	templateName := "template/markdown"
	fsys := createFS(t, []file{{templateName + ".md", `# {{upper .Title}}`}})

	cfg := &ts.Configuration{
		TemplatesFS: &fsys,
		Formats: ts.Formats{
			ts.Text: &ts.FormatOptions{
				FileExtension: "md",
				IsRequired:    true,
			},
		},
	}

	tmpls := ts.NewTemplates(*cfg)
	funcMap := ts.FuncMap{"upper": strings.ToUpper}
	vars := map[string]interface{}{"Title": "Hello from markdown!"}
	tmpl := ts.Must(tmpls.Get(templateName)).Funcs(funcMap)
	expectText := `# HELLO FROM MARKDOWN!`
	var text bytes.Buffer
	if err := tmpl.ExecuteText(&text, vars); err != nil {
		t.Fatal(err)
	}
	gotText := text.String()
	if gotText != expectText {
		t.Errorf("incorrect text template processing; expect %#q; got: %#q", expectText, gotText)
	}
}
