package configer

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
	"io"
)

// decoder is a generic decoder interface.
type decoder interface {
	Decode(v any) error
}

// encoder is a generic encoder interface.
type encoder interface {
	Encode(v any) error
}

// encoding is a generic encoding interface.
// It provides methods to create decoder and encoder.
// Used to abstract the encoding implementation from the Configer.
type encoding interface {
	NewDecoder(r io.Reader) decoder
	NewEncoder(w io.Writer) encoder
}

// # Specific encoding implementations

// ## JSON

type jsonEncoding struct{}

func (e *jsonEncoding) NewDecoder(r io.Reader) decoder {
	return json.NewDecoder(r)
}

func (e *jsonEncoding) NewEncoder(w io.Writer) encoder {
	return json.NewEncoder(w)
}

// JSON encoding for Configer.
// Supported by encoding/json.
var JSON = &jsonEncoding{}

// ## YAML

type yamlEncoding struct{}

func (e *yamlEncoding) NewDecoder(r io.Reader) decoder {
	return yaml.NewDecoder(r)
}

func (e *yamlEncoding) NewEncoder(w io.Writer) encoder {
	return yaml.NewEncoder(w)
}

// YAML encoding for Configer.
// Supported by gopkg.in/yaml.v3.
var YAML = &yamlEncoding{}

// ## TOML

// stdTomlDecoder is a wrapper of toml.Decoder.
// Make it implement the decoder interface.
type stdTomlDecoder struct {
	dec *toml.Decoder
}

// Decode drops the returned Metadata from toml.Decoder.Decode.
func (d *stdTomlDecoder) Decode(v any) error {
	_, err := d.dec.Decode(v)
	return err
}

type tomlEncoding struct{}

func (e *tomlEncoding) NewDecoder(r io.Reader) decoder {
	return &stdTomlDecoder{dec: toml.NewDecoder(r)}
}

func (e *tomlEncoding) NewEncoder(w io.Writer) encoder {
	return toml.NewEncoder(w)
}

// TOML encoding for Configer.
// Supported by github.com/BurntSushi/toml.
var TOML = &tomlEncoding{}
