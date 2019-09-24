package v_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/lyraproj/hierasdk/v"
)

func TestEqual_slice(t *testing.T) {
	expectEqual(t, v.Slice{v.Int(1), v.Int(2)}, v.Slice{v.Int(1), v.Int(2)})
	expectNotEqual(t, v.Slice{v.Int(1), v.Int(2)}, v.Slice{v.Int(1), v.Int(2), v.Int(3)})
	expectNotEqual(t, v.Slice{v.Int(1), v.Int(2)}, v.Slice{v.Int(1), v.Int(3)})
	expectNotEqual(t, v.Slice{v.Int(1), v.Int(2)}, v.Map{`1`: v.Int(3)})
}

func TestEqual_map(t *testing.T) {
	expectEqual(t, v.Map{`one`: v.Int(1), `two`: v.Int(2)}, v.Map{`one`: v.Int(1), `two`: v.Int(2)})
	expectNotEqual(t, v.Map{`one`: v.Int(1), `two`: v.Int(2)}, v.Map{`one`: v.Int(1), `two`: v.Int(2), `three`: v.Int(3)})
	expectNotEqual(t, v.Map{`one`: v.Int(1), `two`: v.Int(2)}, v.Map{`one`: v.Int(1), `two`: v.Int(3)})
	expectNotEqual(t, v.Map{`one`: v.Int(1), `two`: v.Int(2)}, v.Slice{v.Int(1), v.Int(2)})
}

func TestMarshalJSON_bool(t *testing.T) {
	d, err := json.Marshal(v.Bool(true))
	if err != nil {
		t.Fatal(err)
	}
	if string(d) != `true` {
		t.Fatal(`true isn't json 'true'"`)
	}
}

func TestMarshalJSON_map(t *testing.T) {
	m := make(v.Map)
	m[`bool`] = v.Bool(true)
	m[`string`] = v.String(`hello`)
	m[`nil`] = nil
	m[`int`] = v.Int(3)
	m[`float`] = v.Float(3.1)
	m[`map`] = v.Map(map[string]v.Data{
		`a`: v.String(`value of a`),
		`b`: v.String(`value of b`)})
	m[`slice`] = v.Slice{
		v.Bool(false),
		v.Int(1),
		v.Float(2.4),
		nil,
		v.Slice{v.String(`a`), v.String(`b`)}}
	d, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	exp := `{"bool":true,"float":3.1,"int":3,"map":{"a":"value of a","b":"value of b"},"nil":null,"slice":[false,1,2.4,null,["a","b"]],"string":"hello"}`
	if string(d) != exp {
		t.Fatalf(`%s isn't json '%s'`, string(d), exp)
	}
}

func TestUnmarshalJSON_map(t *testing.T) {
	js := `{"bool":true,"float":3.1,"int":3,"map":{"a":"value of a","b":"value of b"},"nil":null,"slice":[false,1,2.4,null,["a","b"]],"string":"hello"}`
	d := v.UnmarshalJSONData([]byte(js))
	m, ok := d.(v.Map)
	if !ok {
		t.Fatal(`expected Map was not produced`)
	}
	d, ok = m[`slice`]
	if !ok {
		t.Fatal(`expected Map does not contain slice`)
	}
	s, ok := d.(v.Slice)
	if !ok {
		t.Fatal(`expected Slice was not produced`)
	}
	d = s[4]
	ds, ok := d.(v.Slice)
	if !ok {
		t.Fatal(`expected nested Slice was not produced`)
	}
	expectEqual(t, v.Slice{v.String(`a`), v.String(`b`)}, ds)
}

func TestMarshalJSON_binary(t *testing.T) {
	b := v.Binary{1, 2, 3}
	d, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	exp := `{"__ptype":"Binary","__pvalue":"AQID"}`
	if string(d) != exp {
		t.Fatalf(`%s isn't json '%s'`, string(d), exp)
	}
}

func TestUnmarshalJSON_binary(t *testing.T) {
	d := v.UnmarshalJSONData([]byte(`{"__ptype":"Binary","__pvalue":"AQID"}`))
	b, ok := d.(v.Binary)
	if !ok {
		t.Fatal(`expected Binary was not produced`)
	}
	expectEqual(t, v.Binary([]byte{1, 2, 3}), b)
	expectPanic(t, `illegal base64 data`, func() { v.UnmarshalJSONData([]byte(`{"__ptype":"Binary","__pvalue":"AQP"}`)) })
}

func TestMarshalJSON_sensitive(t *testing.T) {
	s := v.Sensitive{Data: v.Map{
		`xqz`: v.String(`obfuscated`),
		`sx`:  v.Int(123),
	}}
	d, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	exp := `{"__ptype":"Sensitive","__pvalue":{"sx":123,"xqz":"obfuscated"}}`
	if string(d) != exp {
		t.Fatalf(`%s isn't json '%s'`, string(d), exp)
	}
}

func TestUnmarshalJSON_sensitive(t *testing.T) {
	js := `{"__ptype":"Sensitive","__pvalue":{"sx":123,"xqz":"obfuscated"}}`
	d := v.UnmarshalJSONData([]byte(js))
	b, ok := d.(v.Sensitive)
	if !ok {
		t.Fatal(`expected Sensitive was not produced`)
	}
	sd := v.Sensitive{Data: v.Map{
		`xqz`: v.String(`obfuscated`),
		`sx`:  v.Int(123),
	}}
	expectEqual(t, b, sd)
}

func TestMarshalJSON_timestamp(t *testing.T) {
	now := time.Now()
	ts := now.Format(time.RFC3339Nano)
	s := v.Timestamp(now)
	d, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	exp := `{"__ptype":"Timestamp","__pvalue":"` + ts + `"}`
	if string(d) != exp {
		t.Fatalf(`%s isn't json '%s'`, string(d), exp)
	}
	expectPanic(t, `cannot parse`, func() { v.UnmarshalJSONData([]byte(`{"__ptype":"Timestamp","__pvalue":"bogus"}`)) })
}

func TestMarshalJSON_unknown(t *testing.T) {
	expectPanic(t, `unable to unmarshal`, func() { v.UnmarshalJSONData([]byte(`{"__ptype":"Unknown","__pvalue":"bogus"}`)) })
}

func TestUnmarshalJSON_timestamp(t *testing.T) {
	now := time.Now()
	ts := now.Format(time.RFC3339Nano)
	js := `{"__ptype":"Timestamp","__pvalue":"` + ts + `"}`
	d := v.UnmarshalJSONData([]byte(js))
	b, ok := d.(v.Timestamp)
	if !ok {
		t.Fatal(`expected Timestamp was not produced`)
	}
	expectEqual(t, b, v.Timestamp(now))
}

func TestUnmarshalJSONData_bad(t *testing.T) {
	expectPanic(t, `invalid character`, func() { v.UnmarshalJSONData([]byte(`{1: "one"}`)) })
}

func TestString_binary(t *testing.T) {
	expectStringEqual(t, `Binary("AQID")`, v.Binary{1, 2, 3})
}

func TestString_bool(t *testing.T) {
	expectStringEqual(t, `true`, v.Bool(true))
	expectStringEqual(t, `false`, v.Bool(false))
}

func TestString_float(t *testing.T) {
	expectStringEqual(t, `3.14`, v.Float(3.14))
}

func TestString_int(t *testing.T) {
	expectStringEqual(t, `314`, v.Int(314))
}

func TestString_map(t *testing.T) {
	expectStringEqual(t, `{}`, v.Map{})
	expectStringEqual(t, `{a:"va",b:"vb"}`, v.Map{`a`: v.String(`va`), `b`: v.String(`vb`)})
}

func TestString_sensitive(t *testing.T) {
	expectStringEqual(t, `Sensitive("value redacted")`, v.Sensitive{Data: v.Int(23)})
}

func TestString_slice(t *testing.T) {
	expectStringEqual(t, `[]`, v.Slice{})
	expectStringEqual(t, `["va","vb"]`, v.Slice{v.String(`va`), v.String(`vb`)})
}

func TestString_timestamp(t *testing.T) {
	ts := time.Now()
	expectStringEqual(t, `Timestamp("`+ts.Format(time.RFC3339Nano)+`")`, v.Timestamp(ts))
}

func TestToData_int(t *testing.T) {
	expectEqual(t, v.ToData(5), v.Int(5))
	expectEqual(t, v.ToData(int8(5)), v.Int(5))
	expectEqual(t, v.ToData(int16(5)), v.Int(5))
	expectEqual(t, v.ToData(int32(5)), v.Int(5))
	expectEqual(t, v.ToData(int64(5)), v.Int(5))
	expectEqual(t, v.ToData(uint(5)), v.Int(5))
	expectEqual(t, v.ToData(uint8(5)), v.Int(5))
	expectEqual(t, v.ToData(uint16(5)), v.Int(5))
	expectEqual(t, v.ToData(uint32(5)), v.Int(5))
	expectEqual(t, v.ToData(uint64(5)), v.Int(5))
}

func TestToData_float(t *testing.T) {
	expectEqual(t, v.ToData(float32(5.0)), v.Float(float32(5.0)))
	expectEqual(t, v.ToData(5.0), v.Float(5.0))
}

func TestToData_map_bad(t *testing.T) {
	expectPanic(t, `unable to unmarshal type name`, func() { v.ToData(map[string]int{`__ptype`: 3}) })
}

func TestToData_bad(t *testing.T) {
	expectPanic(t, `unable to create Data from struct`, func() { v.ToData(struct{ A string }{A: `a`}) })
}

func expectEqual(t *testing.T, a, b v.Data) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf(`expected %s, got %s`, a, b)
	}
}

func expectNotEqual(t *testing.T, a, b v.Data) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf(`did not expected %s and %s to be equal`, a, b)
	}
}

func expectStringEqual(t *testing.T, a string, b interface{}) {
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

func expectPanic(t *testing.T, msg string, f func()) {
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
			t.Fatalf(`expected match for %s, got %s`, msg, actual)
		}
	}()
	f()
}
