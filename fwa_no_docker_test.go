//go:build !docker

package fwa

// The build constraint assumes faktory was already set up
func setupFaktory() (func() error, error) {
	return func() error { return nil }, nil
}
