// Package routes provides the Register() function that creates the http.Handler. That
// function is useful when writing tests using the "net/http/httptest" package.
package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/register"
)

func callDataDig(q url.Values, f interface{}) dgo.Value {
	if k := q.Get(`key`); k != `` {
		v, err := vf.UnmarshalJSON([]byte(k))
		if err != nil {
			panic(err)
		}
		if key, ok := v.(dgo.Array); ok {
			return f.(hiera.DataDig)(hiera.NewProviderContext(q), key)
		}
	}
	return nil
}

func callDataHash(q url.Values, f interface{}) dgo.Value {
	return f.(hiera.DataHash)(hiera.NewProviderContext(q))
}

func callLookupKey(q url.Values, f interface{}) dgo.Value {
	if key := q.Get(`key`); key != `` {
		return f.(hiera.LookupKey)(hiera.NewProviderContext(q), key)
	}
	return nil
}

func catch(f func() error) (err error) {
	defer func() {
		switch e := recover().(type) {
		case nil:
		case error:
			err = e
		case string:
			err = errors.New(e)
		default:
			err = fmt.Errorf("error %v", e)
		}
	}()
	err = f()
	return
}

func handleLookup(w http.ResponseWriter, r *http.Request, f func(url.Values, interface{}) dgo.Value, luFunc interface{}) {
	if r.Method != http.MethodGet {
		http.Error(w, ``, http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	err := catch(func() error {
		if r := f(q, luFunc); r != nil {
			return sendData(w, r)
		}
		http.Error(w, `404 value not found`, http.StatusNotFound)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendData(w http.ResponseWriter, d dgo.Value) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(d)
}

// Register create a http.ServeMux and add handlers to it for all lookup functions that has been registered with
// register.DataDig, register.DataHash, and register.LookupKey. The created ServeMux is returned along with a
// Map keyed by function type where each value is a Slice of function names.
func Register() (http.Handler, dgo.Map) {
	if register.Empty() {
		panic(errors.New(`no lookup functions have been registered`))
	}

	router := http.NewServeMux()

	var dataDigNames []dgo.Value
	var dataHashNames []dgo.Value
	var lookupKeyNames []dgo.Value

	register.EachDataDig(func(name string, f hiera.DataDig) {
		dataDigNames = append(dataDigNames, vf.String(name))
		router.HandleFunc(`/data_dig/`+name, func(w http.ResponseWriter, r *http.Request) {
			handleLookup(w, r, callDataDig, f)
		})
	})
	register.EachDataHash(func(name string, f hiera.DataHash) {
		dataHashNames = append(dataHashNames, vf.String(name))
		router.HandleFunc(`/data_hash/`+name, func(w http.ResponseWriter, r *http.Request) {
			handleLookup(w, r, callDataHash, f)
		})
	})
	register.EachLookupKey(func(name string, f hiera.LookupKey) {
		lookupKeyNames = append(lookupKeyNames, vf.String(name))
		router.HandleFunc(`/lookup_key/`+name, func(w http.ResponseWriter, r *http.Request) {
			handleLookup(w, r, callLookupKey, f)
		})
	})
	m := vf.MutableMap(nil)
	if len(dataDigNames) > 0 {
		m.Put(`data_dig`, vf.Array(dataDigNames))
	}
	if len(dataHashNames) > 0 {
		m.Put(`data_hash`, vf.Array(dataHashNames))
	}
	if len(lookupKeyNames) > 0 {
		m.Put(`lookup_key`, vf.Array(lookupKeyNames))
	}
	return router, m
}
