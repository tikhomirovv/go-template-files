package templates

import "io/fs"

type Configuration struct {
	templatesFS *fs.FS
	Formats     Formats
	// globalCssFile  string
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
		templatesFS: templatesFS,
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
