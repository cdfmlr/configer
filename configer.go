package config

import (
	"io"
	"os"
)

// Object is a configuration definition.
// Any
type Object = any

// Configer is a generic configuration reader and writer.
//
// It reads configuration from an io.Reader or file, decodes it with the
// provided encoding, and stores the values into the Config field.
//
// It also comes with Write/WriteToFile methods to reverse the process:
// it encodes the Config object and writes it to an io.Writer or file.
type Configer[T Object] struct {
	Config   *T
	encoding encoding
}

// NewConfiger creates a new Configer.
//
// It binds the configuration it reads to the provided config object.
// encoding is the encoding used to decode and encode the config.
func NewConfiger[T Object](config *T, encoding encoding) *Configer[T] {
	return &Configer[T]{
		Config:   config,
		encoding: encoding,
	}
}

// Read and decode the configuration from the provided io.Reader.
// The result is bound to the c.Config.
func (c *Configer[T]) Read(src io.Reader) error {
	//return yaml.NewDecoder(src).Decode(c.Config)
	return c.encoding.NewDecoder(src).Decode(c.Config)
}

// Write encodes c.Config with c.encoding and writes it to the provided dst.
func (c *Configer[T]) Write(dst io.Writer) error {
	// return yaml.NewEncoder(dst).Encode(c.Config)
	return c.encoding.NewEncoder(dst).Encode(c.Config)
}

// ReadFromFile reads the configuration from the provided file.
// Decodes it with c.encoding and binds the result to c.Config.
func (c *Configer[T]) ReadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.Read(f)
}

// WriteToFile encodes the c.Config with c.encoding and writes it to a file.
func (c *Configer[T]) WriteToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.Write(f)
}
