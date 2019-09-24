package hierasdk

import (
	"context"

	"github.com/lyraproj/hierasdk/v"
)

type (
	// ProviderContext provides utility functions to a provider function
	ProviderContext interface {
		context.Context

		// Option returns the given option or nil if no such option exists
		Option(option string) v.Data

		// reflectToData converts the given value into Data
		ToData(value interface{}) v.Data
	}
)

// A DataDig function looks up a value by a key consisting of several segments. The segments are
// either strings or ints. No other types of segments are allowed.
type DataDig func(ic ProviderContext, key []interface{}) v.Data

// A DataMap function returns a Map that Hiera can use as the source for lookups.
type DataMap func(ic ProviderContext) v.Data

// A LookupKey function returns the value that corresponds to the given key.
type LookupKey func(ic ProviderContext, key string) v.Data

// RegisterDataDig registers a DataDig function under the given name
func RegisterDataDig(name string, f DataDig) {

}

// RegisterDataMap registers a DataMap function under the given name
func RegisterDataMap(name string, f DataMap) {

}

// RegisterLookupKey registers a LookupKey function under the given name
func RegisterLookupKey(name string, f LookupKey) {

}
