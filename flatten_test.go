package gojsonflatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

// deepEquals compares two maps for equality and reports any mismatches.
func deepEquals(t *testing.T, testNumber int, got map[string]interface{}, want map[string]interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%d: mismatch, got: %v wanted: %v", testNumber+1, got, want)
	}
}

func TestFlatten(t *testing.T) {
	cases := []struct {
		test   string
		want   map[string]interface{}
		prefix string
		style  SeparatorStyle
		depth  int
	}{
		// Test case 1
		{
			// JSON input
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			// Expected flattened result
			map[string]interface{}{
				"foo.jim":      "bean",
				"fee":          "bar",
				"n1.alist.0":   "a",
				"n1.alist.1":   "b",
				"n1.alist.2":   "c",
				"n1.alist.3.d": "other",
				"n1.alist.3.e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			// Prefix, SeparatorStyle, and depth
			"",
			DotStyle,
			-1,
		},
		// Test case 2
		{
			// JSON input
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				}
			}`,
			// Expected flattened result
			map[string]interface{}{
				"foo[jim]":        "bean",
				"fee":             "bar",
				"n1[alist][0]":    "a",
				"n1[alist][1]":    "b",
				"n1[alist][2]":    "c",
				"n1[alist][3][d]": "other",
				"n1[alist][3][e]": "another",
			},
			// Prefix, SeparatorStyle, and depth
			"",
			RailsStyle,
			-1,
		},
		// Test case 3
		{
			// JSON input
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			// Expected flattened result
			map[string]interface{}{
				"foo/jim":      "bean",
				"fee":          "bar",
				"n1/alist/0":   "a",
				"n1/alist/1":   "b",
				"n1/alist/2":   "c",
				"n1/alist/3/d": "other",
				"n1/alist/3/e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
		},
		// Test case 4
		{
			// JSON input
			`{ "a": { "b": "c" }, "e": "f" }`,
			// Expected flattened result
			map[string]interface{}{
				"p:a.b": "c",
				"p:e":   "f",
			},
			// Prefix, SeparatorStyle, and depth
			"p:",
			DotStyle,
			-1,
		},
		// Test case 5
		{
			// JSON input
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			// Expected flattened result
			map[string]interface{}{
				"foo_jim":      "bean",
				"fee":          "bar",
				"n1_alist_0":   "a",
				"n1_alist_1":   "b",
				"n1_alist_2":   "c",
				"n1_alist_3_d": "other",
				"n1_alist_3_e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			// Prefix, SeparatorStyle, and depth
			"",
			UnderscoreStyle,
			-1,
		},
		// Test case 6
		{
			// JSON input
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			// Expected flattened result
			map[string]interface{}{
				"flag-foo_jim":      "bean",
				"flag-fee":          "bar",
				"flag-n1_alist_0":   "a",
				"flag-n1_alist_1":   "b",
				"flag-n1_alist_2":   "c",
				"flag-n1_alist_3_d": "other",
				"flag-n1_alist_3_e": "another",
				"flag-number":       1.4567,
				"flag-bool":         true,
			},
			// Prefix, SeparatorStyle, and depth
			"flag-",
			UnderscoreStyle,
			-1,
		},
	}

	for i, test := range cases {
		t.Run(fmt.Sprintf("test: %v", i), func(t *testing.T) {
			var m interface{}
			err := json.Unmarshal([]byte(test.test), &m)
			assert.NoError(t, err)
			got, err := Flatten(m.(map[string]interface{}), test.prefix, test.style, test.depth)
			assert.NoError(t, err)
			deepEquals(t, i, got, test.want)
		})
	}
}

func TestFlattenString(t *testing.T) {
	cases := []struct {
		test   string
		want   string
		prefix string
		style  SeparatorStyle
		depth  int
		err    error
	}{
		// Test case 1
		{
			// JSON input
			`{ "a": "b" }`,
			// Expected flattened result
			`{ "a": "b" }`,
			// Prefix, SeparatorStyle, and depth
			"",
			DotStyle,
			-1,
			nil,
		},
		// Test case 2
		{
			// JSON input
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			// Expected flattened result
			`{ "a.b.c.d": "e", "bool": true, "number": 1.4567 }`,
			// Prefix, SeparatorStyle, and depth
			"",
			DotStyle,
			-1,
			nil,
		},
		// Test case 3
		{
			// JSON input
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			// Expected flattened result
			`{ "a/b/c/d": "e", "bool": true, "number": 1.4567 }`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			nil,
		},
		// Test case 4
		{
			// JSON input
			`{ "a": { "b" : { "c" : { "d" : "e" } } } }`,
			// Expected flattened result
			`{ "a--b--c--d": "e" }`,
			// Prefix, SeparatorStyle, and depth
			"",
			SeparatorStyle{Middle: "--"}, // emdash
			-1,
			nil,
		},
		// Test case 5
		{
			// JSON input
			`{ "a": { "b" : { "c" : { "d" : "e" } } } }`,
			// Expected flattened result
			`{ "a(b)(c)(d)": "e" }`,
			// Prefix, SeparatorStyle, and depth
			"",
			SeparatorStyle{Before: "(", After: ")"}, // paren groupings
			-1,
			nil,
		},
		// Test case 6
		{
			// JSON input with leading whitespace
			`
			  	{ "a": { "b" : { "c" : { "d" : "e" } } } }`,
			// Expected flattened result
			`{ "a(b)(c)(d)": "e" }`,
			// Prefix, SeparatorStyle, and depth
			"",
			SeparatorStyle{Before: "(", After: ")"}, // paren groupings
			-1,
			nil,
		},
		// Test case 7 - Invalid JSON input
		{
			// Invalid JSON input
			`[ "a": { "b": "c" }, "d" ]`,
			// Expected result
			`bogus`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			ErrNotValidJsonInput,
		},
		// Test case 8 - Invalid JSON input
		{
			// Empty JSON input
			``,
			// Expected result
			`bogus`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			ErrNotValidJsonInput,
		},
		// Test case 9 - Invalid JSON input
		{
			// Invalid JSON input (missing quotes)
			`astring`,
			// Expected result
			`bogus`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			ErrNotValidJsonInput,
		},
		// Test case 10 - Invalid JSON input
		{
			// Invalid JSON input (not an object)
			`false`,
			// Expected result
			`bogus`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			ErrNotValidJsonInput,
		},
		// Test case 11 - Invalid JSON input
		{
			// Invalid JSON input (not an object)
			`42`,
			// Expected result
			`bogus`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			ErrNotValidJsonInput,
		},
		// Test case 12 - Invalid JSON input (null)
		{
			// JSON input with null
			`null`,
			// Expected result (error)
			`{}`,
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			-1,
			ErrNotValidJsonInput,
		},
		// Test case 13
		{
			// JSON input
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			// Expected flattened result
			`{ "flag-a.b.c.d": "e", "flag-bool": true, "flag-number": 1.4567 }`,
			// Prefix, SeparatorStyle, and depth
			"flag-",
			DotStyle,
			-1,
			nil,
		},
	}

	nixws := func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}

	for i, test := range cases {
		t.Run(fmt.Sprintf("test: %v", i), func(t *testing.T) {
			got, err := FlattenString(test.test, test.prefix, test.style, test.depth)
			errors.Is(err, test.err)
			if err == nil {
				assert.Equal(t, got, strings.Map(nixws, test.want))
			}
		})
	}
}

func TestFlattenWithVaryingDepth(t *testing.T) {
	cases := []struct {
		test   string
		want   map[string]interface{}
		prefix string
		style  SeparatorStyle
		depth  int
	}{
		// Test data for depth 1
		{
			// JSON input
			`{
                    "foo": {
                        "bar": {
                            "baz": "qux"
                        }
                    },
                    "quux": 789.1
            }`,
			// Expected flattened result
			map[string]interface{}{
				"foo.bar": map[string]interface{}{
					"baz": "qux",
				},
				"quux": 789.1,
			},
			// Prefix, SeparatorStyle, and depth
			"",
			DotStyle,
			1,
		},
		// Test data for depth 2
		{
			// JSON input
			`{
				"a": {
					"b": {
						"c": {
							"d": {
								"e": "f"
							}
						}
					}
				},
		        "g": "h"
		    }`,
			// Expected flattened result
			map[string]interface{}{
				"a_b_c": map[string]interface{}{
					"d": map[string]interface{}{
						"e": "f",
					},
				},
				"g": "h",
			},
			// Prefix, SeparatorStyle, and depth
			"",
			UnderscoreStyle,
			2,
		},

		// Test data for depth 3
		{
			// JSON input
			`{
			"a": {
				"b": {
					"c": {
						"d": {
							"e": "f"
						}
					}
				}
			},
			"g": "h"
		}`,
			// Expected flattened result
			map[string]interface{}{
				"a/b/c/d": map[string]interface{}{
					"e": "f",
				},
				"g": "h",
			},
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			3,
		},

		// Test data for depth 4
		{
			// JSON input
			`{
				"a": {
					"b": {
						"c": {
							"d": {
								"e": 22.2
							}
						}
					}
				},
				"g": "h"
			}`,
			// Expected flattened result
			map[string]interface{}{
				"a/b/c/d/e": 22.2,
				"g":         "h",
			},
			// Prefix, SeparatorStyle, and depth
			"",
			PathStyle,
			4,
		},
		// Test data for depth 5
		{
			// JSON input
			`{
				"a": {
					"b": {
						"c": {
							"d": {
								"e": 22.2
							}
						}
					}
				},
				"g": "h"
			}`,
			// Expected flattened result
			map[string]interface{}{
				"test_a[b][c][d][e]": 22.2,
				"test_g":             "h",
			},
			// Prefix, SeparatorStyle, and depth
			"test_",
			RailsStyle,
			5,
		},
	}

	for i, test := range cases {
		t.Run(fmt.Sprintf("test: %v, depth: %d", i, test.depth), func(t *testing.T) {
			var m interface{}
			err := json.Unmarshal([]byte(test.test), &m)
			assert.NoError(t, err)
			got, err := Flatten(m.(map[string]interface{}), test.prefix, test.style, test.depth)
			assert.NoError(t, err)
			deepEquals(t, i, got, test.want)
		})
	}
}
