package register_test

import (
	"testing"

	expect "github.com/lyraproj/hierasdk/hiera_test"
	"github.com/lyraproj/hierasdk/vf"

	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/register"
)

func TestDataDig(t *testing.T) {
	register.Clean()
	register.DataDig(`l1`, func(ic hiera.ProviderContext, key vf.Slice) vf.Data {
		return nil
	})
	register.DataDig(`l2`, func(ic hiera.ProviderContext, key vf.Slice) vf.Data {
		return nil
	})
	x := ``
	register.EachDataDig(func(n string, _ hiera.DataDig) {
		x += n
	})
	expect.StringEqual(t, `l1l2`, x)
	expect.Panic(t, `already registered`, func() {
		register.DataDig(`l2`, func(ic hiera.ProviderContext, key vf.Slice) vf.Data {
			return nil
		})
	})
}

func TestDataHash(t *testing.T) {
	register.Clean()
	register.DataHash(`l1`, func(ic hiera.ProviderContext) vf.Data {
		return nil
	})
	register.DataHash(`l2`, func(ic hiera.ProviderContext) vf.Data {
		return nil
	})
	x := ``
	register.EachDataHash(func(n string, _ hiera.DataHash) {
		x += n
	})
	expect.StringEqual(t, `l1l2`, x)
}

func TestLookupKey(t *testing.T) {
	register.Clean()
	register.LookupKey(`l1`, func(ic hiera.ProviderContext, key string) vf.Data {
		return nil
	})
	register.LookupKey(`l2`, func(ic hiera.ProviderContext, key string) vf.Data {
		return nil
	})
	x := ``
	register.EachLookupKey(func(n string, _ hiera.LookupKey) {
		x += n
	})
	expect.StringEqual(t, `l1l2`, x)
}

func TestCombo(t *testing.T) {
	register.Clean()
	register.DataDig(`l1`, func(ic hiera.ProviderContext, key vf.Slice) vf.Data {
		return nil
	})
	register.DataHash(`l1`, func(ic hiera.ProviderContext) vf.Data {
		return nil
	})
	register.LookupKey(`l1`, func(ic hiera.ProviderContext, key string) vf.Data {
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
	expect.StringEqual(t, `l1l1l1`, x)
}
