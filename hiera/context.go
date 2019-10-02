package hiera

import (
	"net/url"

	"github.com/lyraproj/hierasdk/vf"
)

type (
	// ProviderContext provides utility functions to a provider function
	ProviderContext interface {
		// Option returns the given option or nil if no such option exists
		Option(option string) vf.Data

		// StringOption returns the option for the given name as a string and true provided that the option is present
		// and is a string. If its missing, or if its found to be something other than a string, this
		// method returns the empty string, false
		StringOption(option string) (string, bool)

		// BoolOption returns the option for the given name as a bool and true provided that the option is present
		// and is a bool. If its missing, or if its found to be something other than a bool, this
		// method returns false, false
		BoolOption(option string) (bool, bool)

		// IntOption returns the option for the given name as an int and true provided that the option is present
		// and is an int. If its missing, or if its found to be something other than an int, this method returns 0, false
		IntOption(option string) (int, bool)

		// FloatOption returns the option for the given name as a float64 and true provided that the option is present
		// and is an float64. If its missing, or if its found to be something other than a float64, this method
		// returns 0.0, false
		FloatOption(option string) (float64, bool)

		// ToData converts the given value into Data
		ToData(value interface{}) vf.Data
	}

	providerContext struct {
		options vf.Map
	}
)

// NewProviderContext creates a context containing the values of the the "options" key in the given url.Values.
func NewProviderContext(q url.Values) ProviderContext {
	var opts vf.Map
	if jo := q.Get(`options`); jo != `` {
		if om, ok := vf.UnmarshalJSONData([]byte(jo)).(vf.Map); ok {
			opts = om
		}
	}
	return &providerContext{options: opts}
}

func (c *providerContext) Option(name string) (d vf.Data) {
	if c.options != nil {
		d = c.options[name]
	}
	return
}

func (c *providerContext) StringOption(name string) (s string, ok bool) {
	var o vf.String
	if o, ok = c.Option(name).(vf.String); ok {
		s = string(o)
	}
	return
}

func (c *providerContext) IntOption(name string) (i int, ok bool) {
	var o vf.Int
	if o, ok = c.Option(name).(vf.Int); ok {
		i = int(o)
	}
	return
}

func (c *providerContext) FloatOption(name string) (f float64, ok bool) {
	var o vf.Float
	if o, ok = c.Option(name).(vf.Float); ok {
		f = float64(o)
	}
	return
}

func (c *providerContext) BoolOption(name string) (b bool, ok bool) {
	var o vf.Bool
	if o, ok = c.Option(name).(vf.Bool); ok {
		b = bool(o)
	}
	return
}

func (c *providerContext) ToData(value interface{}) vf.Data {
	return vf.ToData(value)
}
