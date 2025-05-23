// this file contains functions defined by Grug that can be run from inside a Go template
package main

import "encoding/json"

var grugFuncMap = map[string]any{
	"mkSlice": mkSlice,
	"mkMap":   mkMap,
}

func mkSlice(args ...any) []any {
	// unfortunately I don't know a way to make this work with generics
	// because this needs to be put in the funcMap
	return args
}

// Given a string containing a JSON object, unmarshals that JSON to returns the resulting struct.
func mkMap(obj string) (any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(obj), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
