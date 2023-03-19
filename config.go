package templates

import "io/fs"

type Configuration struct {
	TemplatesFS         *fs.FS
	Formats             Formats
	CommonTemplatesPath *string
}

type Format int

const (
	Html Format = iota
	Text
)

type Formats map[Format]*FormatOptions
type FormatOptions struct {
	FileExtension string
	IsRequired    bool
}

// NewConfiguration creates a default configuration that can be changed
func NewConfiguration(templatesFS *fs.FS) *Configuration {
	return &Configuration{
		TemplatesFS: templatesFS,
		Formats: Formats{
			Html: &FormatOptions{
				FileExtension: "html",
				IsRequired:    true,
			},
			Text: &FormatOptions{
				FileExtension: "txt",
				IsRequired:    false,
			},
		},
	}
}

// SetCommonTemplatesPath sets common templates path
func (c *Configuration) SetCommonTemplatesPath(path string) {
	c.CommonTemplatesPath = &path
}
