package vf_test

import (
	"errors"
	"testing"

	expect "github.com/lyraproj/hierasdk/hiera_test"
	"github.com/lyraproj/hierasdk/vf"
)

func ensureFailed(t *testing.T, f func(t *testing.T)) {
	t.Helper()
	tt := testing.T{}
	f(&tt)
	if !tt.Failed() {
		t.Error(`expected failure did not occur`)
	}
}

func TestExpect(t *testing.T) {
	ensureFailed(t, func(ft *testing.T) {
		expect.Equal(ft, vf.String(`a`), vf.String(`b`))
	})

	ensureFailed(t, func(ft *testing.T) {
		expect.NotEqual(ft, vf.String(`a`), vf.String(`a`))
	})

	ensureFailed(t, func(ft *testing.T) {
		expect.StringEqual(ft, `a`, `b`)
	})

	ensureFailed(t, func(ft *testing.T) {
		expect.StringEqual(ft, `a`, 1)
	})

	ensureFailed(t, func(ft *testing.T) {
		expect.StringEqual(ft, `a`, vf.String(`b`))
	})

	ensureFailed(t, func(ft *testing.T) {
		expect.Panic(ft, `nope`, func() {})
	})

	expect.Panic(t, `oops`, func() { panic(`oops`) })
	expect.Panic(t, `oops`, func() { panic(errors.New(`oops`)) })
	expect.Panic(t, `32`, func() { panic(32) })
}
