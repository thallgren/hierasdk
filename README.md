# Software Development Kit for hiera lookup functions

[![](https://goreportcard.com/badge/github.com/lyraproj/hierasdk)](https://goreportcard.com/report/github.com/lyraproj/hierasdk)
[![](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/lyraproj/hierasdk)
[![](https://github.com/lyraproj/hierasdk/workflows/Hiera%20SDK%20Build/badge.svg)](https://github.com/lyraproj/hierasdk/actions)

This module provides the API to create and run RESTful Hiera lookup plugins. Such a plugin can publish a set of Hiera
lookup functions (i.e. data_dig, data_hash, or lookup_key functions) and make them available using https RESTful
calls.

## How to use
hierasdk is a Go module and is best installed using the command:
```
go get github.com/lyraproj/hierasdk`
```
Provided that Go version >= 1.13 is used, or that the `GO111MODULE=on` is set, this will add the dependency to
the `go.mod` file of the plugin module.

### Example plugin
Skeleton plugin that provides one single `data_hash` function:
```go
package main

import (
  "github.com/lyraproj/hierasdk/hiera"
  "github.com/lyraproj/hierasdk/plugin"
  "github.com/lyraproj/hierasdk/register"
  "github.com/lyraproj/hierasdk/vf"
)

func main() {
  // Register the data_hash function with the global registry
  register.DataHash(`my_data_hash`, myDataHash)

  // Start RESTful service that makes all registered functions available
  plugin.ServeAndExit()
}

func myDataHash(c hiera.ProviderContext) vf.Data {
  var dh map[string]interface{}
  // lookup data hash here and return it
  return vf.ToData(dh)
}
```
