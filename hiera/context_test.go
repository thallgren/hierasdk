package hiera

import (
	"testing"

	expect "github.com/lyraproj/hierasdk/hiera_test"

	"github.com/lyraproj/hierasdk/vf"
)

func TestProviderContext_StringOption(t *testing.T) {
	c := &providerContext{options: vf.Map{`s`: vf.String(`a`), `i`: vf.Int(2)}}
	s, ok := c.StringOption(`s`)
	expect.True(t, ok)
	expect.Equals(t, `a`, s)
	_, ok = c.StringOption(`i`)
	expect.False(t, ok)
}

func TestProviderContext_IntOption(t *testing.T) {
	c := &providerContext{options: vf.Map{`s`: vf.String(`a`), `i`: vf.Int(2)}}
	i, ok := c.IntOption(`i`)
	expect.True(t, ok)
	expect.Equals(t, 2, i)
	_, ok = c.IntOption(`s`)
	expect.False(t, ok)
}

func TestProviderContext_FloatOption(t *testing.T) {
	c := &providerContext{options: vf.Map{`s`: vf.String(`a`), `f`: vf.Float(2)}}
	f, ok := c.FloatOption(`f`)
	expect.True(t, ok)
	expect.Equals(t, 2.0, f)
	_, ok = c.FloatOption(`s`)
	expect.False(t, ok)
}

func TestProviderContext_BoolOption(t *testing.T) {
	c := &providerContext{options: vf.Map{`s`: vf.String(`a`), `b`: vf.Bool(false)}}
	b, ok := c.BoolOption(`b`)
	expect.True(t, ok)
	expect.Equals(t, false, b)
	_, ok = c.BoolOption(`s`)
	expect.False(t, ok)
}
