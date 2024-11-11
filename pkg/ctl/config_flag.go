package ctl

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/berquerant/pneutrinoutil/pkg/reflectx"
	"github.com/spf13/pflag"
)

// ApplyFlagValues sets flag values to this.
func (c *Config) ApplyFlagValues(fs *pflag.FlagSet) error {
	s, err := reflectx.NewFields(*c)
	if err != nil {
		return err
	} // get value of pointer to mutate
	v := reflect.ValueOf(c)

	for _, f := range s {
		name := f.Tag.Get("yaml")
		fv := v.Elem().FieldByName(f.Name)

		switch f.Type.Kind() {
		case reflect.String:
			x, _ := fs.GetString(name)
			fv.SetString(x)
		case reflect.Int:
			x, _ := fs.GetInt(name)
			fv.SetInt(int64(x))
		case reflect.Float32:
			x, _ := fs.GetFloat32(name)
			fv.SetFloat(float64(x))
		default:
			return fmt.Errorf("unexpected Config: %s, %s", name, f.Name)
		}
	}

	return nil
}

// SetFlags sets command-line flags.
// Flag name is from "yaml" struct tag.
// Flag usage is from "usage" struct tag.
// Flag default value is from "default" struct tag.
func (c Config) SetFlags(fs *pflag.FlagSet) error {
	s, err := reflectx.NewFields(c)
	if err != nil {
		return err
	}
	for _, f := range s {
		var (
			name        = f.Tag.Get("yaml")
			usage       = f.Tag.Get("usage")
			dv, dvExist = f.Tag.Lookup("default")
			newErr      = func() error { return fmt.Errorf("unexpected default value of Config: %s, got %s", f.Name, dv) }
		)
		switch f.Type.Kind() {
		case reflect.String:
			var v string
			if dvExist {
				v = dv
			}
			fs.String(name, v, usage)
		case reflect.Int:
			var v int
			if dvExist {
				i, err := strconv.Atoi(dv)
				if err != nil {
					return errors.Join(err, newErr())
				}
				v = i
			}
			fs.Int(name, v, usage)
		case reflect.Float32:
			var v float32
			if dvExist {
				f, err := strconv.ParseFloat(dv, 32)
				if err != nil {
					return errors.Join(err, newErr())
				}
				v = float32(f)
			}
			fs.Float32(name, v, usage)
		default:
			return newErr()
		}
	}

	return nil
}
