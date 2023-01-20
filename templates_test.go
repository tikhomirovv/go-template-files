package templates_test

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"

	templates "github.com/tikhomirovv/go-email-templates"
)

type file struct {
	name     string
	contents string
}

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

func TestGetTemplate(t *testing.T) {

	fsys := createFS(t, []file{
		{"template1/one.html", `<html><head><title>{{.Title}}</title></head></html>`},
		{"template1/one.txt", `# {{.Title}}`},
		{"template2/subdir/two.html", `<html><head><title>{{.Title}}</title></head></html>`},
	})

	cfg := templates.NewConfiguration(&fsys)

	t.Run("HTML:required", func(t *testing.T) {
		cfg.Formats.Html.IsRequired = true
		templates := templates.NewTemplates(*cfg)

		tmpl, err := templates.GetTemplate("template1/no-template", nil)
		testGetTemplateNotExists(t, tmpl, err)

		tmpl, err = templates.GetTemplate("template1/one", nil)
		testGetTemplateExists(t, tmpl, err)
	})

	t.Run("HTML:optional", func(t *testing.T) {
		cfg.Formats.Html.IsRequired = false
		templates := templates.NewTemplates(*cfg)

		tmpl, err := templates.GetTemplate("template1/no-template", nil)
		testGetTemplateExists(t, tmpl, err)
	})

	t.Run("TXT:required", func(t *testing.T) {
		cfg.Formats.Text.IsRequired = true
		templates := templates.NewTemplates(*cfg)

		tmpl, err := templates.GetTemplate("template2/subdir/two", nil)
		testGetTemplateNotExists(t, tmpl, err)

		tmpl, err = templates.GetTemplate("template1/one", nil)
		testGetTemplateExists(t, tmpl, err)
	})

	t.Run("TXT:optional", func(t *testing.T) {
		cfg.Formats.Text.IsRequired = false
		templates := templates.NewTemplates(*cfg)

		tmpl, err := templates.GetTemplate("template2/subdir/two", nil)
		testGetTemplateExists(t, tmpl, err)
	})

	t.Run("Content:check", func(t *testing.T) {
		templates := templates.NewTemplates(*cfg)
		tmpl, _ := templates.GetTemplate("template1/one", map[string]interface{}{"Title": "Hello from tests!"})
		expectHtml := `<html><head><title>Hello from tests!</title></head></html>`
		expectText := `# Hello from tests!`

		if tmpl.Html != expectHtml {
			t.Errorf("incorrect html template processing; expect %#q; got: %#q", expectHtml, tmpl.Html)
		}
		if tmpl.Text != expectText {
			t.Errorf("incorrect text template processing; expect %#q; got: %#q", expectText, tmpl.Text)
		}
	})
}

func testGetTemplateExists(t *testing.T, tmpl *templates.Template, err error) {
	if err != nil {
		t.Error(err)
	}

	if tmpl == nil {
		t.Error("the template must be found")
	}
}

func testGetTemplateNotExists(t *testing.T, tmpl *templates.Template, err error) {
	if _, ok := err.(templates.FileNotFoundError); err != nil && !ok {
		t.Errorf("the error must be of type %#q ", "FileNotFoundError")
	}

	if tmpl != nil {
		t.Error("the template must not be found")
	}
}
