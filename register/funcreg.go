package register

import (
	"fmt"
	"sort"
	"sync"

	"github.com/lyraproj/hierasdk/hiera"
)

type (
	funcReg struct {
		lock       sync.RWMutex
		dataDigs   map[string]interface{}
		dataHashes map[string]interface{}
		lookupKeys map[string]interface{}
	}
)

var global = funcReg{}

// EachDataDig calls the given actor once with each registered DataDig function
func (r *funcReg) EachDataDig(actor func(name string, f hiera.DataDig)) {
	r.sortedEach(r.dataDigs, func(n string, f interface{}) { actor(n, f.(hiera.DataDig)) })
}

// EachDataHash calls the given actor once with each registered DataHash function
func (r *funcReg) EachDataHash(actor func(name string, f hiera.DataHash)) {
	r.sortedEach(r.dataHashes, func(n string, f interface{}) { actor(n, f.(hiera.DataHash)) })
}

// EachLookupKey calls the given actor once with each registered LookupKey function
func (r *funcReg) EachLookupKey(actor func(name string, f hiera.LookupKey)) {
	r.sortedEach(r.lookupKeys, func(n string, f interface{}) { actor(n, f.(hiera.LookupKey)) })
}

// Empty returns true if no functions have been registered
func (r *funcReg) Empty() bool {
	r.lock.RLock()
	empty := len(r.dataDigs)+len(r.dataHashes)+len(r.lookupKeys) == 0
	r.lock.RUnlock()
	return empty
}

// DataDig registers a DataDig function under the given name
func (r *funcReg) DataDig(name string, f hiera.DataDig) {
	r.register(&r.dataDigs, `data_dig`, name, f)
}

// DataHash registers a DataHash function under the given name
func (r *funcReg) DataHash(name string, f hiera.DataHash) {
	r.register(&r.dataHashes, `data_hash`, name, f)
}

// LookupKey registers a LookupKey function under the given name
func (r *funcReg) LookupKey(name string, f hiera.LookupKey) {
	r.register(&r.lookupKeys, `lookup_key`, name, f)
}

func (r *funcReg) sortedEach(m map[string]interface{}, f func(string, interface{})) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	ks := make([]string, len(m))
	i := 0
	for k := range m {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	for _, n := range ks {
		f(n, m[n])
	}
}

func (r *funcReg) register(mp *map[string]interface{}, tp, name string, f interface{}) {
	r.lock.Lock()
	m := *mp
	if m == nil {
		m = make(map[string]interface{})
		*mp = m
	}
	if _, ok := m[name]; ok {
		r.lock.Unlock()
		panic(fmt.Errorf(`%s function '%s' is already registered`, tp, name))
	}
	m[name] = f
	r.lock.Unlock()
}

// Clean removes any prior registrations. Should only be used in tests
func Clean() {
	global = funcReg{}
}

// DataDig registers a DataDig function under the given name with the global registry
func DataDig(name string, f hiera.DataDig) {
	global.DataDig(name, f)
}

// DataHash registers a DataHash function under the given name with the global registry
func DataHash(name string, f hiera.DataHash) {
	global.DataHash(name, f)
}

// LookupKey registers a LookupKey function under the given name with the global registry
func LookupKey(name string, f hiera.LookupKey) {
	global.LookupKey(name, f)
}

// EachDataDig calls the given actor once with each registered DataDig function in the global registry
func EachDataDig(actor func(name string, f hiera.DataDig)) {
	global.EachDataDig(actor)
}

// EachDataHash calls the given actor once with each registered DataHash function in the global registry
func EachDataHash(actor func(name string, f hiera.DataHash)) {
	global.EachDataHash(actor)
}

// EachLookupKey calls the given actor once with each registered LookupKey function in the global registry
func EachLookupKey(actor func(name string, f hiera.LookupKey)) {
	global.EachLookupKey(actor)
}

// Empty returns true if no functions have been registered with the global registry
func Empty() bool {
	return global.Empty()
}
