# TypeID Go
### A golang implementation of [TypeIDs](https://github.com/jetpack-io/typeid)
![License: Apache 2.0](https://img.shields.io/github/license/jetpack-io/typeid-go) [![Go Reference](https://pkg.go.dev/badge/go.jetpack.io/typeid.svg)](https://pkg.go.dev/go.jetpack.io/typeid)

TypeIDs are a modern, **type-safe**, globally unique identifier based on the upcoming
UUIDv7 standard. They provide a ton of nice properties that make them a great choice
as the primary identifiers for your data in a database, APIs, and distributed systems.
Read more about TypeIDs in their [spec](https://github.com/jetpack-io/typeid).

This particular implementation provides a go library for generating and parsing TypeIDs.

## Installation

To add this library as a dependency in your go module, run:

```bash
go get go.jetpack.io/typeid
```

## Usage
This library provides both a statically typed and a dynamically typed version of TypeIDs.

The statically typed version lives under the `typed` package. It makes it possible for
the go compiler itself to enforce type safety.

To use it, first define your TypeID types:

```go
import (
  typeid "go.jetpack.io/typeid/typed"
)

type userPrefix struct{}
func (userPrefix) Type() string { return "user" }
type UserID struct { typeid.TypeID[userPrefix] }
```

And now use those types to generate TypeIDs:

```go
import (
  typeid "go.jetpack.io/typeid/typed"
)

func example() {
  tid := typeid.New[UserID]()
  fmt.Println(tid)
}
```

If you don't want static types, you can use the dynamic version instead:
  
```go
import (
  "go.jetpack.io/typeid/typeid"
)

func example() {
  tid := typeid.New("user")
  fmt.Println(tid)
}
```

For the full documentation, see this package's [godoc](https://pkg.go.dev/go.jetpack.io/typeid).