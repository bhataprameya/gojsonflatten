// This Package Flattens Json with the specified key separators up to a specified depth
package gojsonflatten

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

// SeparatorStyle defines the style of keys when flattening nested structures.
// It provides options for adding prefixes, middle separators, and suffixes.
type SeparatorStyle struct {
	Before string // Prepend to key
	Middle string // Add between keys
	After  string // Append to key
}

// Default SeparatorStyles
var (
	DotStyle        = SeparatorStyle{Middle: "."}
	PathStyle       = SeparatorStyle{Middle: "/"}
	RailsStyle      = SeparatorStyle{Before: "[", After: "]"}
	UnderscoreStyle = SeparatorStyle{Middle: "_"}
)

var (
	ErrNotValidInput     = errors.New("not a valid input: map or slice")
	ErrNotValidJsonInput = errors.New("not a valid input, must be a map")
	isJsonMap            = regexp.MustCompile(`^\s*\{`)
)

// Flatten generates a flat map from a nested map with a specified depth.
func Flatten(nested map[string]interface{}, prefix string, style SeparatorStyle, depth int) (map[string]interface{}, error) {
	return flattenInternal(nested, prefix, style, depth, false)
}

// FlattenString generates a flat JSON map from a nested JSON string with a specified depth.
func FlattenString(nestedString, prefix string, style SeparatorStyle, depth int) (string, error) {
	return flattenStringInternal(nestedString, prefix, style, depth, false)
}

// FlattenNoArray generates a flat map from a nested map with a specified depth, preserving arrays as strings.
func FlattenNoArray(nested map[string]interface{}, prefix string, style SeparatorStyle, depth int) (map[string]interface{}, error) {
	return flattenInternal(nested, prefix, style, depth, true)
}

// FlattenStringNoArray generates a flat JSON map from a nested JSON string with a specified depth, preserving arrays as strings.
func FlattenStringNoArray(nestedString, prefix string, style SeparatorStyle, depth int) (string, error) {
	return flattenStringInternal(nestedString, prefix, style, depth, true)
}

// flattenInternal generates a flat map from a nested map with a specified depth, optionally preserving arrays as strings.
func flattenInternal(nested map[string]interface{}, prefix string, style SeparatorStyle, depth int, preserveArray bool) (map[string]interface{}, error) {
	if depth == 0 {
		return nested, nil
	} else if depth > 0 {
		depth++
	}

	flatmap := make(map[string]interface{})
	err := flatten(true, flatmap, nested, prefix, style, depth, preserveArray)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}

// flattenStringInternal generates a flat JSON map from a nested JSON string with a specified depth, optionally preserving arrays as strings.
func flattenStringInternal(nestedString, prefix string, style SeparatorStyle, depth int, preserveArray bool) (string, error) {
	if !isJsonMap.MatchString(nestedString) {
		return "", ErrNotValidJsonInput
	}

	var nested map[string]interface{}
	err := json.Unmarshal([]byte(nestedString), &nested)
	if err != nil {
		return "", err
	}

	flatmap, err := flattenInternal(nested, prefix, style, depth, preserveArray)
	if err != nil {
		return "", err
	}

	flatBytes, err := json.Marshal(&flatmap)
	if err != nil {
		return "", err
	}

	return string(flatBytes), nil
}

// flatten recursively processes nested structures and flattens them.
func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, style SeparatorStyle, depth int, keepArrays bool) error {
	if depth == 0 {
		// If the desired depth is reached, add the prefix and nested value to the flat map.
		flatMap[prefix] = nested
		return nil
	}

	// Assign function is used to process and assign values to the flat map.
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			// If the value is a nested map or slice, continue flattening recursively.
			if err := flatten(false, flatMap, v, newKey, style, depth-1, keepArrays); err != nil {
				return err
			}
		default:
			// For scalar values, directly add them to the flat map.
			flatMap[newKey] = v
		}

		return nil
	}

	switch nested := nested.(type) {
	case map[string]interface{}:
		for k, v := range nested {
			newKey := enkey(top, prefix, k, style)
			// Process and assign the key-value pair.
			assign(newKey, v)
		}
	case []interface{}:
		if !keepArrays {
			for i, v := range nested {
				newKey := enkey(top, prefix, strconv.Itoa(i), style)
				// Process and assign the index-value pair.
				assign(newKey, v)
			}
		} else {
			flatMap[prefix] = nested
		}
	default:
		return ErrNotValidInput
	}

	return nil
}

// enkey combines the prefix, subKey, and SeparatorStyle to form a new key.
func enkey(top bool, prefix, subKey string, style SeparatorStyle) string {
	key := prefix

	if top {
		// If it's the top level, directly add the subKey.
		key += subKey
	} else {
		// For nested levels, use the specified SeparatorStyle.
		key += style.Before + style.Middle + subKey + style.After
	}

	return key
}
