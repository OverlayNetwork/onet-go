package onet

import "github.com/libs4go/sdi4go"

// Config .
type Config struct {
	injector sdi4go.Injector
}

// NewConfig .
func NewConfig() *Config {
	return &Config{
		injector: sdi4go.New(),
	}
}

// Bind bind option object with name
func (builder *Config) Bind(name string, obj interface{}) error {
	return builder.injector.Bind(name, sdi4go.Singleton(obj))
}

// Get get name option
func (builder *Config) Get(name string, objptr interface{}) bool {
	err := builder.injector.Create(name, objptr)

	if err != nil {
		return false
	}

	return true
}

// Option config overlay network start option
type Option func(*Config) error
