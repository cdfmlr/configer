package configer

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
	"io"
)

// Decoder is a generic Decoder interface.
type Decoder interface {
	Decode(v any) error
}

// Encoder is a generic Encoder interface.
type Encoder interface {
	Encode(v any) error
}

// Encoding is a generic Encoding interface.
// It provides methods to create Decoder and Encoder.
// Used to abstract the Encoding implementation from the Configer.
type Encoding interface {
	NewDecoder(r io.Reader) Decoder
	NewEncoder(w io.Writer) Encoder
}

// # Specific Encoding implementations

// ## JSON

type jsonEncoding struct{}

func (e *jsonEncoding) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func (e *jsonEncoding) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

// JSON Encoding for Configer.
// Supported by Encoding/json.
var JSON = &jsonEncoding{}

// ## YAML

type yamlEncoding struct{}

func (e *yamlEncoding) NewDecoder(r io.Reader) Decoder {
	return yaml.NewDecoder(r)
}

func (e *yamlEncoding) NewEncoder(w io.Writer) Encoder {
	return yaml.NewEncoder(w)
}

// YAML Encoding for Configer.
// Supported by gopkg.in/yaml.v3.
var YAML = &yamlEncoding{}

// ## TOML

// stdTomlDecoder is a wrapper of toml.Decoder.
// Make it implement the Decoder interface.
type stdTomlDecoder struct {
	dec *toml.Decoder
}

// Decode drops the returned Metadata from toml.Decoder.Decode.
func (d *stdTomlDecoder) Decode(v any) error {
	_, err := d.dec.Decode(v)
	return err
}

type tomlEncoding struct{}

func (e *tomlEncoding) NewDecoder(r io.Reader) Decoder {
	return &stdTomlDecoder{dec: toml.NewDecoder(r)}
}

func (e *tomlEncoding) NewEncoder(w io.Writer) Encoder {
	return toml.NewEncoder(w)
}

// TOML Encoding for Configer.
// Supported by github.com/BurntSushi/toml.
var TOML = &tomlEncoding{}
