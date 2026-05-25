# variant

[![Go Reference](https://pkg.go.dev/badge/github.com/mababaNiubi/variant.svg)](https://pkg.go.dev/github.com/mababaNiubi/variant)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[中文](README_CN.md)

A generic, dynamically-typed value container for Go that turns data into a weakly-typed variant. **Variant** is a discriminated union that can hold any of eight runtime types — Empty, Bool, Int, UInt, Float, String, List, Map — with JSON/MessagePack serialization, mixed-type arithmetic, and reflection-based decoding from native Go values.

Designed for data transmission, collection, and telemetry pipelines where values arrive as loosely-typed strings and need efficient numeric computation and storage.

## Installation

```bash
go get github.com/mababaNiubi/variant
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/mababaNiubi/variant"
)

func main() {
    // Smart constructor — auto-detects type
    v := variant.New("3.14")
    fmt.Println(v.Type())        // TypeFloat64
    fmt.Println(v.AsFloat64())   // 3.14 <nil>

    // JSON parsing
    v2 := variant.New(`{"name":"Alice","age":30}`)
    name, _ := v2.MapGet("name")
    fmt.Println(name.AsString()) // "Alice"

    // Arithmetic
    a := variant.NewInt(10)
    b := variant.NewFloat64(3.5)
    result, _ := a.Increase(b)
    fmt.Println(result.AsFloat64()) // 13.5 <nil>

    // JSON serialization
    data, _ := v2.MarshalJSON()
    fmt.Println(string(data))    // {"age":30,"name":"Alice"}
}
```

## Constructors

```go
v := variant.New(x)          // Smart factory — accepts any Go value
v := variant.NewEmpty()      // TypeEmpty
v := variant.NewBool(true)   // TypeBool
v := variant.NewInt(42)      // TypeInt64
v := variant.NewInt64(42)    // TypeInt64
v := variant.NewUInt64(42)   // TypeUInt64
v := variant.NewValue(3.14)  // Auto int/float detection
v := variant.NewFloat64(3.14)// TypeFloat64
v := variant.NewString("hi") // TypeString
v := variant.NewValueList([]variant.Variant{...})   // TypeList
v := variant.NewValueMap(map[string]variant.Variant{...}) // TypeMap
```

`New()` accepts: `bool`, `string`, `float64`, `float32`, `int`, `int8`–`int64`, `uint`–`uint64`, `[]byte`, `Variant`, `*Variant`, `[]Variant`, `map[string]Variant`, and any struct (via reflection).

For string inputs, `New()` tries to parse as: int → float → empty → JSON → plain string.


## Type Conversion

```go
b, err := v.AsBool()       // → bool
i, err := v.AsInt()        // → int
i64, err := v.AsInt64()    // → int64
u64, err := v.AsUInt64()   // → uint64
f32, err := v.AsFloat32()  // → float32
f64, err := v.AsFloat64()  // → float64
s := v.AsString()          // → string (never fails)
s := v.String()            // → string (implements fmt.Stringer)
iface := v.AsInterface()   // → recursive conversion to plain Go types
```

Conversions handle cross-type coercion with overflow protection.

## Equality & Comparison

```go
v.IsEqual(other)                              // deep comparison
ok, err := v.CompareNumberBySymbol(r, ">=")   // numeric comparison
ok := v.Comparable(r)                         // numeric comparsion with lexicographic fallback
```

Comparison symbols: `=`, `!=`, `>`, `<`, `>=`, `<=`.

## Arithmetic

```go
result, err := a.Increase(b)  // a + b
result, err := a.Reduce(b)    // a - b
result, err := a.Multiple(b)  // a * b
result, err := a.Divide(b)    // a / b
result := v.Decimal(2)        // round to 2 decimal places
```

Mixed-type rules:

- Int + Float → Float
- String + anything → String concatenation
- List + List → merged list
- List + scalar → append
- Map + Map → merged map (overwrites duplicate keys)

## Container Operations

```go
// Polymorphic (accept int index or string key)
v.Add(value)               // append to list or concat to string
v.Set(indexOrKey, value)   // set by index (list) or key (map)
val, ok := v.Get(idxOrKey) // get by index or key
v.Remove(idxOrKey)         // remove by index or key

// Type-specific
val, ok := v.ListGet(0)
v.ListSet(0, value)
val, ok := v.MapGet("key")
v.MapSet("key", value)     // initializes empty variant to TypeMap

// Iteration
v.Range(func(key string, val Variant) bool { ... })
v.RangeByIndex(func(idx int, val Variant) bool { ... })
```

Return `false` from the callback to stop iteration early.

## JSON Serialization

```go
// Marshal
data, err := v.MarshalJSON()
json.Marshal(v)  // via json.Marshaler interface

// Unmarshal
var v variant.Variant
v.UnmarshalJSON([]byte(`{"key":"value"}`))

// Standalone
v, err := variant.UnmarshalJSON([]byte(`[1,2,3]`))
```

## Binary Serialization (MessagePack)

```go
// With format marker (for WAL / storage)
data, err := v.MarshalBinary()
v, n, err := variant.UnmarshalBinary(data)
ok := variant.IsBinaryFormat(data)

// Without marker (for embedding in larger msgpack payloads)
data, err := v.MarshalMsgpack()
v.UnmarshalMsgpack(data)
```

The binary format prepends a `0x01` byte to distinguish from legacy JSON.

## Decoding Go Structs

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
v := variant.New(Person{Name: "Alice", Age: 30})
// v is TypeMap with keys "name" and "age" (from json tags)
```

Uses reflection. Struct fields use `json` tags for key names, falling back to the Go field name.

## Performance Design

- **Inline numeric storage**: Integers, booleans, and float64 values are stored in an `int64` field via `unsafe.Pointer` bit-casting — no heap allocation for the common numeric path.
- **Zero-allocation msgpack encoding**: `encodeVariant` writes directly to the encoder without intermediate allocations.
- **Lazy string parsing**: Numeric strings are not parsed until a typed conversion is requested.

## License

MIT — see [LICENSE](LICENSE).
