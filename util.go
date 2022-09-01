package json

import (
	"fmt"
)

// ErrNoTable indicates that a chart does not have a matching table.
type ErrNoTable struct {
	Key string
}

func (e ErrNoTable) Error() string { return fmt.Sprintf("%q is not a table", e.Key) }

// ErrNoValue indicates that Values does not contain a key with a value
type ErrNoValue struct {
	Key string
}

func (e ErrNoValue) Error() string { return fmt.Sprintf("%q is not a value", e.Key) }

// Merges source and into map, preferring values from the source map ( values -> into)
func mergeValues(into map[string]interface{}, values map[string]interface{}) map[string]interface{} {
	for k, v := range values {
		// If the key doesn't exist already, then just set the key to that value
		if _, ok := into[k]; !ok {
			into[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			into[k] = v
			continue
		}
		intoMap, isMap := into[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			into[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		into[k] = mergeValues(intoMap, nextMap)
	}
	return into
}

func clonePath(path []string) []string {
	ret := make([]string, len(path))
	copy(ret, path)
	return ret
}
