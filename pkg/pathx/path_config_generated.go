// Code generated by "goconfig -field Mode os.FileMode|Truncate bool -option -output path_config_generated.go"; DO NOT EDIT.

package pathx

import "os"

type ConfigItem[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func (s *ConfigItem[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *ConfigItem[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}
func (s *ConfigItem[T]) Default() T {
	return s.defaultValue
}
func (s *ConfigItem[T]) IsModified() bool {
	return s.modified
}
func NewConfigItem[T any](defaultValue T) *ConfigItem[T] {
	return &ConfigItem[T]{
		defaultValue: defaultValue,
	}
}

type Config struct {
	Mode     *ConfigItem[os.FileMode]
	Truncate *ConfigItem[bool]
}
type ConfigBuilder struct {
	mode     os.FileMode
	truncate bool
}

func (s *ConfigBuilder) Mode(v os.FileMode) *ConfigBuilder {
	s.mode = v
	return s
}
func (s *ConfigBuilder) Truncate(v bool) *ConfigBuilder {
	s.truncate = v
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		Mode:     NewConfigItem(s.mode),
		Truncate: NewConfigItem(s.truncate),
	}
}

func NewConfigBuilder() *ConfigBuilder { return &ConfigBuilder{} }
func (s *Config) Apply(opt ...ConfigOption) {
	for _, x := range opt {
		x(s)
	}
}

type ConfigOption func(*Config)

func WithMode(v os.FileMode) ConfigOption {
	return func(c *Config) {
		c.Mode.Set(v)
	}
}
func WithTruncate(v bool) ConfigOption {
	return func(c *Config) {
		c.Truncate.Set(v)
	}
}