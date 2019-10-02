// Package vf contains everything necessary to ensure type safe JSON serialization of implementors of
// the Data interface. The set of implementations is fixed and hard wired into the ToData function.
package vf

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	// Data represents a data structure that can be serialized using the json.Marshaler
	Data interface {
		fmt.Stringer
		json.Marshaler

		// Equals returns the result of comparing this instance with the given value
		Equals(d Data) bool
	}

	// Binary is the Data type for []byte
	Binary []byte

	// Bool is the Data type for bool
	Bool bool

	// Float is the Data type for float64
	Float float64

	// Int is the Data type for int
	Int int

	// Map is the Data type for map[string]Data
	Map map[string]Data

	// Sensitive is the Data type for sensitive data
	Sensitive struct {
		Data
	}

	// String is the Data type for string
	String string

	// Slice is the Data type for []Data
	Slice []Data

	// Timestamp is the Data type for time.Time
	Timestamp time.Time
)

const typeKey = `__ptype`
const valueKey = `__pvalue`
const sensitiveName = `Sensitive`
const binaryName = `Binary`
const timestampName = `Timestamp`

var rTypeKey = reflect.ValueOf(typeKey)
var rValueKey = reflect.ValueOf(valueKey)

// Equals returns the result of comparing this instance with the given value
func (d Binary) Equals(o Data) bool {
	if cb, ok := o.(Binary); ok {
		return bytes.Equal(d, cb)
	}
	return false
}

// MarshalJSON will output {"__ptype":"Binary","__pvalue":<quoted strict base64 of contained bytes>}
func (d Binary) MarshalJSON() ([]byte, error) {
	return marshalTypeMap(binaryName, base64.StdEncoding.Strict().EncodeToString(d))
}

// String returns 'Binary(<quoted strict base64 of contained bytes>)'
func (d Binary) String() string {
	return fmt.Sprintf(`Binary(%q)`, base64.StdEncoding.Strict().EncodeToString(d))
}

// Equals returns the result of comparing this instance with the given value
func (d Bool) Equals(o Data) bool {
	return d == o
}

// MarshalJSON creates the json encoding for this value
func (d Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(d))
}

// String returns the 'true' or 'false'
func (d Bool) String() string {
	s := `false`
	if d {
		s = `true`
	}
	return s
}

// Equals returns the result of comparing this instance with the given value
func (d Float) Equals(o Data) bool {
	return d == o
}

// MarshalJSON creates the json encoding for this value
func (d Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(d))
}

// String returns the 'g' format for the float
func (d Float) String() string {
	return strconv.FormatFloat(float64(d), 'g', -1, 64)
}

// Equals returns the result of comparing this instance with the given value
func (d Int) Equals(o Data) bool {
	return d == o
}

// MarshalJSON creates the json encoding for this value
func (d Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(d))
}

// String returns the the a string containing a decimal integer
func (d Int) String() string {
	return strconv.Itoa(int(d))
}

// Equals returns the result of comparing this instance with the given value
func (d Map) Equals(o Data) bool {
	if od, ok := o.(Map); ok && len(d) == len(od) {
		for k, av := range d {
			if bv, ok := od[k]; ok && av.Equals(bv) {
				continue
			}
			return false
		}
		return true
	}
	return false
}

// MarshalJSON creates the json encoding for this value
func (d Map) MarshalJSON() ([]byte, error) {
	return json.Marshal((map[string]Data)(d))
}

// String returns '{' <key>:<value> [',' <key>:<value> ...] '}' sorted in alphabetical order of the keys
func (d Map) String() string {
	l := len(d)
	if l == 0 {
		return `{}`
	}
	ks := make([]string, l)
	i := 0
	for k := range d {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	b := strings.Builder{}
	dl := '{'
	for i = range ks {
		k := ks[i]
		_, _ = b.WriteRune(dl)
		dl = ','
		_, _ = b.WriteString(k)
		_, _ = b.WriteRune(':')
		_, _ = b.WriteString(d[k].String())
	}
	_, _ = b.WriteRune('}')
	return b.String()
}

// Equals returns the result of comparing this instance with the given value
func (d Sensitive) Equals(o Data) bool {
	if od, ok := o.(Sensitive); ok {
		return d.Unwrap().Equals(od.Unwrap())
	}
	return false
}

// MarshalJSON will output {"__ptype":"Sensitive","__pvalue":<json of wrapped data>}
func (d Sensitive) MarshalJSON() ([]byte, error) {
	return marshalTypeMap(sensitiveName, d.Data)
}

// String returns 'Sensitive("value redacted")'
func (d Sensitive) String() string {
	return `Sensitive("value redacted")`
}

// Unwrap returns the wrapped data
func (d Sensitive) Unwrap() Data {
	return d.Data
}

// Equals returns the result of comparing this instance with the given value
func (d Slice) Equals(o Data) bool {
	if od, ok := o.(Slice); ok && len(d) == len(od) {
		for i := range d {
			if !d[i].Equals(od[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// MarshalJSON creates the json encoding for this value
func (d Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(([]Data)(d))
}

// String returns '[' <value> [',' <value> ...] ']'
func (d Slice) String() string {
	b := strings.Builder{}
	dl := '['
	for i := range d {
		_, _ = b.WriteRune(dl)
		dl = ','
		_, _ = b.WriteString(d[i].String())
	}
	if dl == '[' {
		_, _ = b.WriteString(`[]`)
	} else {
		_, _ = b.WriteRune(']')
	}
	return b.String()
}

// Equals returns the result of comparing this instance with the given value
func (d String) Equals(o Data) bool {
	return d == o
}

// MarshalJSON creates the json encoding for this value
func (d String) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// String returns this string quoted
func (d String) String() string {
	return strconv.Quote(string(d))
}

// Equals returns the result of comparing this instance with the given value
func (d Timestamp) Equals(o Data) bool {
	if od, ok := o.(Timestamp); ok {
		return time.Time(d).Equal(time.Time(od))
	}
	return false
}

// MarshalJSON will output {"__ptype":"Timestamp","__pvalue":<quoted RFC3339Nano representation of timestamp>}
func (d Timestamp) MarshalJSON() ([]byte, error) {
	return marshalTypeMap(timestampName, time.Time(d))
}

// String returns 'Timestamp(<quoted RFC3339Nano representation of timestamp>)
func (d Timestamp) String() string {
	return fmt.Sprintf(`Timestamp(%q)`, time.Time(d).Format(time.RFC3339Nano))
}

// UnmarshalJSONData transform the given json into a Data value
func UnmarshalJSONData(j []byte) Data {
	var v interface{}
	enc := json.NewDecoder(bytes.NewReader(j))
	enc.UseNumber()
	err := enc.Decode(&v)
	if err != nil {
		panic(err)
	}
	return ToData(v)
}

// ToBool returns the bool that the given Data represents. The method panics unless the Data is an Bool
func ToBool(d Data) bool {
	return bool(d.(Bool))
}

// ToString returns the string that the given Data represents. The method panics unless the Data is a String
func ToString(d Data) string {
	return string(d.(String))
}

// ToInt returns the int that the given Data represents. The method panics unless the Data is an Int
func ToInt(d Data) int {
	return int(d.(Int))
}

// ToFloat returns the float64 that the given Data represents. The method panics unless the Data is an Float
func ToFloat(d Data) float64 {
	return float64(d.(Float))
}

// ToData converts the given value to a Data value
func ToData(v interface{}) Data {
	if d, ok := v.(Data); ok {
		return d
	}
	return reflectToData(reflect.ValueOf(v))
}

func reflectToData(rv reflect.Value) (d Data) {
	k := rv.Kind()
	switch {
	case reflect.Invalid == k:
		d = nil
	case reflect.Slice == k:
		d = sliceToData(rv)
	case reflect.String == k:
		d = stringToData(rv.Interface())
	case reflect.Map == k:
		d = mapToData(rv)
	case reflect.Interface == k:
		d = reflectToData(rv.Elem())
	case reflect.Bool == k:
		d = Bool(rv.Bool())
	case reflect.Int <= k && k <= reflect.Int64:
		d = Int(rv.Int())
	case reflect.Uint <= k && k <= reflect.Uint64:
		d = Int(rv.Uint())
	case reflect.Float32 == k || reflect.Float64 == k:
		d = Float(rv.Float())
	default:
		panic(fmt.Errorf(`unable to create Data from %#v, kind %d`, rv, rv.Kind()))
	}
	return
}

func stringToData(v interface{}) (d Data) {
	switch v := v.(type) {
	case string:
		d = String(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			d = Int(i)
		} else {
			f, _ := v.Float64()
			d = Float(f)
		}
	default:
		panic(fmt.Errorf(`unable to create Data from %T %#v`, v, v))
	}
	return
}

func sliceToData(rv reflect.Value) Data {
	l := rv.Len()
	r := make(Slice, l)
	for i := 0; i < l; i++ {
		r[i] = reflectToData(rv.Index(i))
	}
	return r
}

func mapToData(rv reflect.Value) Data {
	if tn := rv.MapIndex(rTypeKey); tn.IsValid() {
		dn := reflectToData(tn)
		if s, ok := dn.(String); ok {
			return richToData(string(s), rv)
		}
		panic(fmt.Errorf(`unable to unmarshal type name '%#v'`, tn))
	}
	ks := rv.MapKeys()
	r := make(Map, len(ks))
	for i := range ks {
		kr := ks[i]
		k := reflectToData(kr)
		r[string(k.(String))] = reflectToData(rv.MapIndex(kr))
	}
	return r
}

func richToData(tn string, rv reflect.Value) Data {
	if pv := rv.MapIndex(rValueKey); pv.IsValid() {
		pd := reflectToData(pv)
		switch tn {
		case binaryName:
			if s, ok := pd.(String); ok {
				bs, err := base64.StdEncoding.Strict().DecodeString(string(s))
				if err != nil {
					panic(err)
				}
				return Binary(bs)
			}
		case sensitiveName:
			return Sensitive{reflectToData(pv)}
		case timestampName:
			if s, ok := pd.(String); ok {
				t, err := time.Parse(time.RFC3339Nano, string(s))
				if err != nil {
					panic(err)
				}
				return Timestamp(t)
			}
		}
	}
	panic(fmt.Errorf(`unable to unmarshal '%#v'`, rv))
}

func marshalTypeMap(tn string, v interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		typeKey:  tn,
		valueKey: v,
	})
}
