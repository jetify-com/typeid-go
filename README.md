# TypeID Go

### A golang implementation of [TypeIDs](https://github.com/jetify-com/typeid)

![License: Apache 2.0](https://img.shields.io/github/license/jetpack-io/typeid-go) [![Go Reference](https://pkg.go.dev/badge/go.jetpack.io/typeid.svg)](https://pkg.go.dev/go.jetpack.io/typeid)

TypeIDs are a modern, **type-safe**, globally unique identifier based on the upcoming
UUIDv7 standard. They provide a ton of nice properties that make them a great choice
as the primary identifiers for your data in a database, APIs, and distributed systems.
Read more about TypeIDs in their [spec](https://github.com/jetify-com/typeid).

This particular implementation provides a go library for generating and parsing TypeIDs.

## Installation

To add this library as a dependency in your go module, run:

```bash
go get go.jetpack.io/typeid
```

## Usage

This library provides a go implementation of TypeID that allows you
to define your own custom id types for added compile-time safety.

If you don't need compile-time safety, you can use the provided `typeid.AnyID` directly:

```go
import (
  "go.jetpack.io/typeid"
)

func example() {
  tid, _ := typeid.WithPrefix("user")
  fmt.Println(tid)
}
```

If you want compile-time safety, define your own custom types with two steps:

1. Define a struct the implements the method `Prefix`. Prefix should return the
   string that should be used as the prefix for your custom type.
2. Define you own id type, by embedding `typeid.TypeID[CustomPrefix]`

For example to define a UserID with prefix `user`:

```go
import (
  "go.jetpack.io/typeid"
)

// Define the prefix:
type UserPrefix struct {}
func (UserPrefix) Prefix() string { return "user" }

// Define UserID:
type UserID struct {
	typeid.TypeID[UserPrefix]
}
```

Now you can use the UserID type to generate new ids:

```go
import (
  "go.jetpack.io/typeid"
)

func example() {
  tid, _ := typeid.New[UserID]()
  fmt.Println(tid)
}
```

For the full documentation, see this package's [godoc](https://pkg.go.dev/go.jetpack.io/typeid).
