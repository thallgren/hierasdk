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

		// ToData converts the given value into Data
		ToData(value interface{}) vf.Data
	}

	providerContext struct {
		options vf.Map
	}
)

func NewProviderContext(q url.Values) ProviderContext {
	var opts vf.Map
	if jo := q.Get(`options`); jo != `` {
		if om, ok := vf.UnmarshalJSONData([]byte(jo)).(vf.Map); ok {
			opts = om
		}
	}
	return &providerContext{options: opts}
}

func (c *providerContext) Option(option string) vf.Data {
	if c.options != nil {
		return c.options[option]
	}
	return nil
}

func (c *providerContext) ToData(value interface{}) vf.Data {
	return vf.ToData(value)
}
