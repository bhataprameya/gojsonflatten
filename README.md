# Json & Map Flatten Library

[![Go Reference](https://pkg.go.dev/badge/github.com/jeremywohl/flatten.svg)](https://pkg.go.dev/github.com/bhataprameya/gojsonflatten)

Flatten is a Go library that simplifies the process of converting arbitrarily nested JSON structures into flat, one-dimensional maps. It offers flexibility in how you want to format the flattened keys, supporting default styles such as dotted, path-like, Rails, or underscores, or even allowing you to define your custom style. Whether you're working with JSON strings or Go structures, Flatten can handle it all, and you can specify the depth to flatten, from flattening any depth (-1) to a specific level.

## Features

- Flatten nested JSON structures into one-dimensional maps.
- Specify the flattening depth or flatten any depth.
- Support various key styles: dotted, path-like, Rails, or underscores.
- Define custom key styles to suit your needs.
- Works with both JSON strings and Go structures.
- Ability to flatten only maps and not arrays.

## Installation

To use the Flatten library in your Go project, you can simply install it using `go get`:

```bash
go get github.com/bhataprameya/gojsonflatten
```

## Usage

### Flatten JSON Strings

You can flatten JSON strings using Flatten as shown below:

```GO
import (
    "github.com/bhataprameya/gojsonflatten"
)

nested := `{
  "one": {
    "two": [
      "2a",
      "2b"
    ]
  },
  "side": "value"
}`

flat, err := gojsonflatten.FlattenString(nested, "", flatten.DotStyle, -1)

// Output:
// {
//   "one.two.0": "2a",
//   "one.two.1": "2b",
//   "side": "value"
// }
```

### Flatten Go Maps

You can also flatten Go maps directly like this:

```GO

import (
    "github.com/bhataprameya/gojsonflatten"
)

nested := map[string]interface{}{
   "a": "b",
   "c": map[string]interface{}{
       "d": "e",
       "f": "g",
   },
   "z": 1.4567,
}

flat, err := gojsonflatten.Flatten(nested, "", flatten.RailsStyle, 3)

// Output:
// {
//   "a":    "b",
//   "c[d]": "e",
//   "c[f]": "g",
//   "z":    1.4567,
// }
```

### Flatten Go Maps without Arrays

You can also flatten Go maps while preserving arrays as it is like this::

```GO

import (
    "github.com/bhataprameya/gojsonflatten"
)

nested := map[string]interface{}{
   "a": "b",
   "c": map[string]interface{}{
       "d": "e",
       "f": "g",
   },
   "z": []interface{}{"one", "two"},
}

flat, err := gojsonflatten.FlattenNoArray(nested, "", flatten.RailsStyle, 3)

// Output:
// {
//   "a":    "b",
//   "c[d]": "e",
//   "c[f]": "g",
//   "z":    []interface{}{"one", "two"}
// }
```

### Flatten JSON Strings without Arrays

You can flatten JSON strings while preserving arrays using FlattenNoArray:

```GO
import (
    "github.com/bhataprameya/gojsonflatten"
)

nested := `{
  "one": {
    "two": [
      "2a",
      "2b"
    ]
  },
  "side": "value"
}`

flat, err := gojsonflatten.FlattenStringNoArray(nested, "", flatten.DotStyle, -1)

// Output:
// {
//   "one.two": ["2a", "2b"],
//   "side": "value"
// }

```

## Custom Key Style

You can even define a custom key style for flattening:

```GO
import (
    "github.com/bhataprameya/gojsonflatten"
)

doubleDash := gojsonflatten.SeparatorStyle{Middle: "--"}
flat, err := gojsonflatten.FlattenString(nested, "", doubleDash, 5)

// Output:
// {
//   "one--two--0": "2a",
//   "one--two--1": "2b",
//   "side": "value"
// }
```

## Prefix Keys

You can even define a custom key style for flattening:

```GO
import (
    "github.com/bhataprameya/gojsonflatten"
)

doubleDash := gojsonflatten.SeparatorStyle{Middle: "--"}
flat, err := gojsonflatten.FlattenString(nested, "hello_", doubleDash, 5)

// Output:
// {
//   "hello_one--two--0": "2a",
//   "hello_one--two--1": "2b",
//   "hello_side": "value"
// }
```

## Contributing

If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request. Your contributions are greatly appreciated!

## License

This library is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Acknowledgments

This library is a fork of `https://github.com/jeremywohl/flatten` and has been updated and modified to add additional features

#### Happy flattening! 🚀
