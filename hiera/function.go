package hiera

import (
	"context"

	"github.com/lyraproj/hierasdk/vf"
)

type (
	// ProviderContext provides utility functions to a provider function
	ProviderContext interface {
		context.Context

		// Option returns the given option or nil if no such option exists
		Option(option string) vf.Data

		// ToData converts the given value into Data
		ToData(value interface{}) vf.Data
	}

	// DataDig is a Hiera 'data_dig' function looks up a value by a key consisting of several segments.
	// The segments are either strings or ints. No other types of segments are allowed.
	DataDig func(ic ProviderContext, key vf.Slice) vf.Data

	// DataHash is a Hiera 'data_hash' function returns a Map that Hiera can use as the source for
	// lookups.
	DataHash func(ic ProviderContext) vf.Data

	// LookupKey is a Hiera 'lookup_key' function returns the value that corresponds to the given key.
	LookupKey func(ic ProviderContext, key string) vf.Data
)
