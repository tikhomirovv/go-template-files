package templates

import "io/fs"

type Format string

type formatOptions struct {
	FileExtension string
	IsRequired    bool
}

type formats struct {
	Html formatOptions
	Text formatOptions
}

type configuration struct {
	templatesFS *fs.FS
	Formats     formats
	// globalCssFile  string
	// htmlIsRequired bool
	// textIsRequired bool
}

func NewConfiguration(templatesFS *fs.FS) *configuration {
	return &configuration{
		templatesFS: templatesFS,
		Formats: formats{
			Html: formatOptions{
				FileExtension: "html",
				IsRequired:    true,
			},
			Text: formatOptions{
				FileExtension: "txt",
				IsRequired:    false,
			},
		},
	}
}

func (fo *formatOptions) SetExtension(ext string) {
	fo.FileExtension = ext
}

func (fo *formatOptions) SetIsRequired(ir bool) {
	fo.IsRequired = ir
}
