package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/register"
	"github.com/lyraproj/hierasdk/vf"
)

func callDataDig(q url.Values, f interface{}) vf.Data {
	if k := q.Get(`key`); k != `` {
		if key, ok := vf.UnmarshalJSONData([]byte(k)).(vf.Slice); ok {
			return f.(hiera.DataDig)(hiera.NewProviderContext(q), key)
		}
	}
	return nil
}

func callDataHash(q url.Values, f interface{}) vf.Data {
	return f.(hiera.DataHash)(hiera.NewProviderContext(q))
}

func callLookupKey(q url.Values, f interface{}) vf.Data {
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

func handleLookup(w http.ResponseWriter, r *http.Request, f func(url.Values, interface{}) vf.Data, luFunc interface{}) {
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

func sendData(w http.ResponseWriter, d vf.Data) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(d)
}

// Register create a http.ServeMux and add handlers to it for all lookup functions that has been registered with
// register.DataDig, register.DataHash, and register.LookupKey. The created ServeMux is returned along with a
// Map keyed by function type where each value is a Slice of function names.
func Register() (http.Handler, vf.Map) {
	if register.Empty() {
		panic(errors.New(`no lookup functions have been registered`))
	}

	router := http.NewServeMux()

	var dataDigNames vf.Slice
	var dataHashNames vf.Slice
	var lookupKeyNames vf.Slice

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
	m := make(vf.Map)
	if len(dataDigNames) > 0 {
		m[`data_dig`] = dataDigNames
	}
	if len(dataHashNames) > 0 {
		m[`data_hash`] = dataHashNames
	}
	if len(lookupKeyNames) > 0 {
		m[`lookup_key`] = lookupKeyNames
	}
	return router, m
}
