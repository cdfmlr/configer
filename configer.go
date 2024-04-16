package config

import (
	"io"
	"os"
)

type Object = any

type Configer[T Object] struct {
	Config   *T
	encoding encoding
}

func NewConfiger[T Object](config *T, encoding encoding) *Configer[T] {
	return &Configer[T]{
		Config:   config,
		encoding: encoding,
	}
}

func (c *Configer[T]) Read(src io.Reader) error {
	//return yaml.NewDecoder(src).Decode(c.Config)
	return c.encoding.NewDecoder(src).Decode(c.Config)
}

func (c *Configer[T]) Write(dst io.Writer) error {
	// return yaml.NewEncoder(dst).Encode(c.Config)
	return c.encoding.NewEncoder(dst).Encode(c.Config)
}

func (c *Configer[T]) ReadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.Read(f)
}

func (c *Configer[T]) WriteToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.Write(f)
}
