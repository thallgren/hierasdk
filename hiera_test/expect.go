package expect

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/lyraproj/hierasdk/vf"
)

// Equals expects a.Equals(b) to be true
func Equals(t *testing.T, a, b interface{}) {
	t.Helper()
	ad := vf.ToData(a)
	bd := vf.ToData(b)
	if !ad.Equals(bd) {
		t.Errorf(`expected %s, got %s`, ad, bd)
	}
}

// NotEqual expects a.Equals(b) to be false
func NotEqual(t *testing.T, a, b vf.Data) {
	t.Helper()
	ad := vf.ToData(a)
	bd := vf.ToData(b)
	if ad.Equals(bd) {
		t.Errorf(`did not expected %s and %s to be equal`, ad, bd)
	}
}

// True expects b to be true
func True(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Errorf(`expected true`)
	}
}

// False expects b to be false
func False(t *testing.T, b bool) {
	t.Helper()
	if b {
		t.Errorf(`expected false`)
	}
}

// StringEqual expects that the string a is equal to b where b is either a string
// or an implementor of fmt.Stringer
func StringEqual(t *testing.T, a string, b interface{}) {
	t.Helper()
	bs, ok := b.(string)
	if !ok {
		if s, ok := b.(fmt.Stringer); ok {
			bs = s.String()
		} else {
			bs = `not a string`
		}
	}
	if a != bs {
		t.Errorf(`expected %q, got %q`, a, bs)
	}
}

// Panic expects that the given function panics with an error or string that
// matches the given msg pattern
func Panic(t *testing.T, msg string, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		var actual string
		if r == nil {
			actual = `no error`
		} else if er, ek := r.(error); ek {
			actual = er.Error()
		} else if er, ek := r.(string); ek {
			actual = er
		} else {
			actual = fmt.Sprintf("%#v", r)
		}
		if !regexp.MustCompile(msg).MatchString(actual) {
			t.Errorf(`expected match for %s, got %s`, msg, actual)
		}
	}()
	f()
}
