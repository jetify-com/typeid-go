# TypeID Go

### A golang implementation of [TypeIDs](https://github.com/jetify-com/typeid)

![License: Apache 2.0](https://img.shields.io/github/license/jetify-com/typeid-go) [![Go Reference](https://pkg.go.dev/badge/go.jetify.com/typeid.svg)](https://pkg.go.dev/go.jetify.com/typeid)

TypeIDs are a modern, **type-safe**, globally unique identifier based on the upcoming
UUIDv7 standard. They provide a ton of nice properties that make them a great choice
as the primary identifiers for your data in a database, APIs, and distributed systems.
Read more about TypeIDs in their [spec](https://github.com/jetify-com/typeid).

This particular implementation provides a go library for generating and parsing TypeIDs.

## Installation

To add this library as a dependency in your go module, run:

```bash
go get go.jetify.com/typeid
```

## Usage

This library provides a go implementation of TypeID:

```go
import (
  "go.jetify.com/typeid"
)

func example() {
  // Generate a new TypeID with a prefix (panics on invalid prefix)
  tid := typeid.MustGenerate("user")
  fmt.Println(tid)
  
  // Generate a new TypeID without a prefix
  tid = typeid.MustGenerate("")
  fmt.Println(tid)
  
  // Generate with error handling
  tid, err := typeid.Generate("user")
  if err != nil {
    log.Fatal(err)
  }
  
  // Parse an existing TypeID
  tid, _ = typeid.Parse("user_00041061050r3gg28a1c60t3gf")
  fmt.Println(tid)
  
  // Convert from UUID
  tid, _ = typeid.FromUUID("user", "018e5f71-6f04-7c5c-8123-456789abcdef")
  fmt.Println(tid)
}
```

For the full documentation, see this package's [godoc](https://pkg.go.dev/go.jetify.com/typeid).
