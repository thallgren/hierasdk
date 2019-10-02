package vf_test

import (
	"encoding/json"
	"testing"
	"time"

	expect "github.com/lyraproj/hierasdk/hiera_test"
	"github.com/lyraproj/hierasdk/vf"
)

func TestEqual_binary(t *testing.T) {
	expect.Equals(t, vf.Binary(`hello`), vf.Binary(`hello`))
	expect.NotEqual(t, vf.Binary(`hello`), vf.Binary(`good bye`))
	expect.NotEqual(t, vf.Binary(`hello`), vf.String(`hello`))
}

func TestEqual_bool(t *testing.T) {
	expect.Equals(t, vf.Bool(true), vf.Bool(true))
	expect.NotEqual(t, vf.Bool(true), vf.Bool(false))
	expect.NotEqual(t, vf.Bool(true), vf.Int(1))
}

func TestEqual_slice(t *testing.T) {
	expect.Equals(t, vf.Slice{vf.Int(1), vf.Int(2)}, vf.Slice{vf.Int(1), vf.Int(2)})
	expect.NotEqual(t, vf.Slice{vf.Int(1), vf.Int(2)}, vf.Slice{vf.Int(1), vf.Int(2), vf.Int(3)})
	expect.NotEqual(t, vf.Slice{vf.Int(1), vf.Int(2)}, vf.Slice{vf.Int(1), vf.Int(3)})
	expect.NotEqual(t, vf.Slice{vf.Int(1), vf.Int(2)}, vf.Map{`1`: vf.Int(3)})
}

func TestEqual_map(t *testing.T) {
	expect.Equals(t, vf.Map{`one`: vf.Int(1), `two`: vf.Int(2)}, vf.Map{`one`: vf.Int(1), `two`: vf.Int(2)})
	expect.NotEqual(t, vf.Map{`one`: vf.Int(1), `two`: vf.Int(2)}, vf.Map{`one`: vf.Int(1), `two`: vf.Int(2), `three`: vf.Int(3)})
	expect.NotEqual(t, vf.Map{`one`: vf.Int(1), `two`: vf.Int(2)}, vf.Map{`one`: vf.Int(1), `two`: vf.Int(3)})
	expect.NotEqual(t, vf.Map{`one`: vf.Int(1), `two`: vf.Int(2)}, vf.Slice{vf.Int(1), vf.Int(2)})
}

func TestEqual_sensitive(t *testing.T) {
	expect.Equals(t, vf.Sensitive{Data: vf.Int(5)}, vf.Sensitive{Data: vf.Int(5)})
	expect.NotEqual(t, vf.Sensitive{Data: vf.Int(5)}, vf.Sensitive{Data: vf.Int(4)})
	expect.NotEqual(t, vf.Sensitive{Data: vf.Int(5)}, vf.Int(5))
}

func TestEqual_timestamp(t *testing.T) {
	now := time.Now()
	expect.Equals(t, vf.Timestamp(now), vf.Timestamp(now))
	expect.NotEqual(t, vf.Timestamp(now), vf.Timestamp(now.Add(1)))
	expect.NotEqual(t, vf.Timestamp(now), vf.Int(now.UnixNano()))
}

func TestMarshalJSON_bool(t *testing.T) {
	d, err := json.Marshal(vf.Bool(true))
	if err != nil {
		t.Fatal(err)
	}
	if string(d) != `true` {
		t.Fatal(`true isn't json 'true'"`)
	}
}

func TestMarshalJSON_map(t *testing.T) {
	m := make(vf.Map)
	m[`bool`] = vf.Bool(true)
	m[`string`] = vf.String(`hello`)
	m[`nil`] = nil
	m[`int`] = vf.Int(3)
	m[`float`] = vf.Float(3.1)
	m[`map`] = vf.Map(map[string]vf.Data{
		`a`: vf.String(`value of a`),
		`b`: vf.String(`value of b`)})
	m[`slice`] = vf.Slice{
		vf.Bool(false),
		vf.Int(1),
		vf.Float(2.4),
		nil,
		vf.Slice{vf.String(`a`), vf.String(`b`)}}
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
	d := vf.UnmarshalJSONData([]byte(js))
	m, ok := d.(vf.Map)
	if !ok {
		t.Fatal(`expected Map was not produced`)
	}
	d, ok = m[`slice`]
	if !ok {
		t.Fatal(`expected Map does not contain slice`)
	}
	s, ok := d.(vf.Slice)
	if !ok {
		t.Fatal(`expected Slice was not produced`)
	}
	d = s[4]
	ds, ok := d.(vf.Slice)
	if !ok {
		t.Fatal(`expected nested Slice was not produced`)
	}
	expect.Equals(t, vf.Slice{vf.String(`a`), vf.String(`b`)}, ds)
}

func TestMarshalJSON_binary(t *testing.T) {
	b := vf.Binary{1, 2, 3}
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
	d := vf.UnmarshalJSONData([]byte(`{"__ptype":"Binary","__pvalue":"AQID"}`))
	b, ok := d.(vf.Binary)
	if !ok {
		t.Fatal(`expected Binary was not produced`)
	}
	expect.Equals(t, vf.Binary([]byte{1, 2, 3}), b)
	expect.Panic(t, `illegal base64 data`, func() { vf.UnmarshalJSONData([]byte(`{"__ptype":"Binary","__pvalue":"AQP"}`)) })
}

func TestMarshalJSON_sensitive(t *testing.T) {
	s := vf.Sensitive{Data: vf.Map{
		`xqz`: vf.String(`obfuscated`),
		`sx`:  vf.Int(123),
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
	d := vf.UnmarshalJSONData([]byte(js))
	b, ok := d.(vf.Sensitive)
	if !ok {
		t.Fatal(`expected Sensitive was not produced`)
	}
	sd := vf.Sensitive{Data: vf.Map{
		`xqz`: vf.String(`obfuscated`),
		`sx`:  vf.Int(123),
	}}
	expect.Equals(t, b, sd)
}

func TestMarshalJSON_timestamp(t *testing.T) {
	now := time.Now()
	ts := now.Format(time.RFC3339Nano)
	s := vf.Timestamp(now)
	d, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	exp := `{"__ptype":"Timestamp","__pvalue":"` + ts + `"}`
	if string(d) != exp {
		t.Fatalf(`%s isn't json '%s'`, string(d), exp)
	}
	expect.Panic(t, `cannot parse`, func() { vf.UnmarshalJSONData([]byte(`{"__ptype":"Timestamp","__pvalue":"bogus"}`)) })
}

func TestMarshalJSON_unknown(t *testing.T) {
	expect.Panic(t, `unable to unmarshal`, func() { vf.UnmarshalJSONData([]byte(`{"__ptype":"Unknown","__pvalue":"bogus"}`)) })
}

func TestUnmarshalJSON_timestamp(t *testing.T) {
	now := time.Now()
	ts := now.Format(time.RFC3339Nano)
	js := `{"__ptype":"Timestamp","__pvalue":"` + ts + `"}`
	d := vf.UnmarshalJSONData([]byte(js))
	b, ok := d.(vf.Timestamp)
	if !ok {
		t.Fatal(`expected Timestamp was not produced`)
	}
	expect.Equals(t, b, vf.Timestamp(now))
}

func TestUnmarshalJSONData_bad(t *testing.T) {
	expect.Panic(t, `invalid character`, func() { vf.UnmarshalJSONData([]byte(`{1: "one"}`)) })
}

func TestString_binary(t *testing.T) {
	expect.StringEqual(t, `Binary("AQID")`, vf.Binary{1, 2, 3})
}

func TestString_bool(t *testing.T) {
	expect.StringEqual(t, `true`, vf.Bool(true))
	expect.StringEqual(t, `false`, vf.Bool(false))
}

func TestString_float(t *testing.T) {
	expect.StringEqual(t, `3.14`, vf.Float(3.14))
}

func TestString_int(t *testing.T) {
	expect.StringEqual(t, `314`, vf.Int(314))
}

func TestString_map(t *testing.T) {
	expect.StringEqual(t, `{}`, vf.Map{})
	expect.StringEqual(t, `{a:"va",b:"vb"}`, vf.Map{`a`: vf.String(`va`), `b`: vf.String(`vb`)})
}

func TestString_sensitive(t *testing.T) {
	expect.StringEqual(t, `Sensitive("value redacted")`, vf.Sensitive{Data: vf.Int(23)})
}

func TestString_slice(t *testing.T) {
	expect.StringEqual(t, `[]`, vf.Slice{})
	expect.StringEqual(t, `["va","vb"]`, vf.Slice{vf.String(`va`), vf.String(`vb`)})
}

func TestString_timestamp(t *testing.T) {
	ts := time.Now()
	expect.StringEqual(t, `Timestamp("`+ts.Format(time.RFC3339Nano)+`")`, vf.Timestamp(ts))
}

func TestToData_int(t *testing.T) {
	expect.Equals(t, vf.ToData(5), vf.Int(5))
	expect.Equals(t, vf.ToData(int8(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(int16(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(int32(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(int64(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(uint(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(uint8(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(uint16(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(uint32(5)), vf.Int(5))
	expect.Equals(t, vf.ToData(uint64(5)), vf.Int(5))
}

func TestToData_data(t *testing.T) {
	s := vf.String(`hello`)
	expect.Equals(t, s, vf.ToData(s))
}

func TestToData_float(t *testing.T) {
	expect.Equals(t, vf.ToData(float32(5.0)), vf.Float(float32(5.0)))
	expect.Equals(t, vf.ToData(5.0), vf.Float(5.0))
}

func TestToData_map_bad(t *testing.T) {
	expect.Panic(t, `unable to unmarshal type name`, func() { vf.ToData(map[string]int{`__ptype`: 3}) })
}

func TestToData_bad(t *testing.T) {
	type stuff string
	expect.Panic(t, `unable to create Data from struct`, func() { vf.ToData(struct{ A string }{A: `a`}) })
	expect.Panic(t, `unable to create Data from vf_test\.stuff`, func() { vf.ToData(stuff(`a`)) })
}
