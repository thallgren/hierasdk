package hiera

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/vf"
)

func TestProviderContextFromMap(t *testing.T) {
	pc := ProviderContextFromMap(nil)
	require.Equal(t, vf.Map(), pc.OptionsMap())

	pc = ProviderContextFromMap(vf.MutableMap(`a`, vf.MutableValues(1, 2, 3)))
	require.True(t, pc.Option(`a`).(dgo.Array).Frozen())
}

func TestProviderContext_StringOption(t *testing.T) {
	c := &providerContext{options: vf.Map(`s`, `a`, `i`, 2)}
	s, ok := c.StringOption(`s`)
	require.True(t, ok)
	require.Equal(t, `a`, s)
	_, ok = c.StringOption(`i`)
	require.False(t, ok)
}

func TestProviderContext_IntOption(t *testing.T) {
	c := &providerContext{options: vf.Map(`s`, `a`, `i`, 2)}
	i, ok := c.IntOption(`i`)
	require.True(t, ok)
	require.Equal(t, 2, i)
	_, ok = c.IntOption(`s`)
	require.False(t, ok)
}

func TestProviderContext_FloatOption(t *testing.T) {
	c := &providerContext{options: vf.Map(`s`, `a`, `f`, 2.0)}
	f, ok := c.FloatOption(`f`)
	require.True(t, ok)
	require.Equal(t, 2.0, f)
	_, ok = c.FloatOption(`s`)
	require.False(t, ok)
}

func TestProviderContext_BoolOption(t *testing.T) {
	c := &providerContext{options: vf.Map(`s`, `a`, `b`, false)}
	b, ok := c.BoolOption(`b`)
	require.True(t, ok)
	require.Equal(t, false, b)
	_, ok = c.BoolOption(`s`)
	require.False(t, ok)
}
