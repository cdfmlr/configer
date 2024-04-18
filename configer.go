package configer

import (
	"io"
	"os"
)

// Configuration is a configuration definition.
//
// Configer binds the configuration it reads to a Configuration object.
// Or it serializes the Configuration object to write it to a file.
//
// Any type can be used as a Configuration as long as it is serializable.
// It is recommended to use a struct type to represent the configuration.
type Configuration = any

// Configer is a generic configuration reader and writer.
//
// It reads configuration from an io.Reader or file, decodes it with the
// provided encoding, and stores the values into the Config field.
//
// It also comes with Write/WriteToFile methods to reverse the process:
// it encodes the Config object and writes it to an io.Writer or file.
//
// The type parameter T is the configuration type.
// Actually, It is not required introducing generics here, but if I am not
// misunderstanding, a type parameter is helpful to keep the Config field type
// without abstracting it to an interface.
// (It's an impulsive decision made at the time of 1.18 release. Keeping it
// for now, as it is at least harmless.)
type Configer[T Configuration] struct {
	// Config is a pointer to the configuration object.
	Config *T
	// encoding to use for marshaling/unmarshalling.
	encoding encoding
}

// New creates a new Configer.
//
// It binds the configuration it reads to the provided config object.
// The encoding parameter specifies the encoding to use for marshaling/unmarshalling.
func New[T Configuration](config *T, encoding encoding) *Configer[T] {
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
