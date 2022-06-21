package templatewrapper

import "io"

// Functionality that this library exposes.
type Api interface {
	Parse() error
	Write(wr io.Writer, name string, data interface{}) error
}
