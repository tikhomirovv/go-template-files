package templates

import "fmt"

type FileNotFoundError string

// Error returns non-empty string if there was an error.
func (filename FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found: %#q", string(filename))
}
