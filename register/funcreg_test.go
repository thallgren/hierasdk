package register_test

import (
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/register"
)

func TestDataDig(t *testing.T) {
	register.Clean()
	register.DataDig(`l1`, func(ic hiera.ProviderContext, key dgo.Array) dgo.Value {
		return nil
	})
	register.DataDig(`l2`, func(ic hiera.ProviderContext, key dgo.Array) dgo.Value {
		return nil
	})
	x := ``
	register.EachDataDig(func(n string, _ hiera.DataDig) {
		x += n
	})
	require.Equal(t, `l1l2`, x)
	require.Panic(t, func() {
		register.DataDig(`l2`, func(ic hiera.ProviderContext, key dgo.Array) dgo.Value {
			return nil
		})
	}, `already registered`)
}

func TestDataHash(t *testing.T) {
	register.Clean()
	register.DataHash(`l1`, func(ic hiera.ProviderContext) dgo.Map {
		return nil
	})
	register.DataHash(`l2`, func(ic hiera.ProviderContext) dgo.Map {
		return nil
	})
	x := ``
	register.EachDataHash(func(n string, _ hiera.DataHash) {
		x += n
	})
	require.Equal(t, `l1l2`, x)
}

func TestLookupKey(t *testing.T) {
	register.Clean()
	register.LookupKey(`l1`, func(ic hiera.ProviderContext, key string) dgo.Value {
		return nil
	})
	register.LookupKey(`l2`, func(ic hiera.ProviderContext, key string) dgo.Value {
		return nil
	})
	x := ``
	register.EachLookupKey(func(n string, _ hiera.LookupKey) {
		x += n
	})
	require.Equal(t, `l1l2`, x)
}

func TestCombo(t *testing.T) {
	register.Clean()
	register.DataDig(`l1`, func(ic hiera.ProviderContext, key dgo.Array) dgo.Value {
		return nil
	})
	register.DataHash(`l1`, func(ic hiera.ProviderContext) dgo.Map {
		return nil
	})
	register.LookupKey(`l1`, func(ic hiera.ProviderContext, key string) dgo.Value {
		return nil
	})
	x := ``
	register.EachDataDig(func(n string, _ hiera.DataDig) {
		x += n
	})
	register.EachDataHash(func(n string, _ hiera.DataHash) {
		x += n
	})
	register.EachLookupKey(func(n string, _ hiera.LookupKey) {
		x += n
	})
	require.Equal(t, `l1l1l1`, x)
}
