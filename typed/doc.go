// Package typed implements a statically checked version of TypeIDs
//
// To use it, define your own ID types that implement the IDType interface:
//
//	type userPrefix struct{}
//	func (userPrefix) Type() string { return "user" }
//	type UserID struct {
//		typeid.TypeID[userPrefix]
//	}
//
//	type accountPrefix struct{}
//	func (accountPrefix) Type() string { return "account" }
//	type AccountID struct {
//		typeid.TypeID[accountPrefix]
//	}
//
// And now you can use your IDTypes via generics. For example, to create a
// new ID of type user:
//
//	import (
//		typeid "go.jetpack.io/typeid/typed"
//	)
//
//	user_id, _ := typeid.New[UserID]()
//
// Because this implementation uses generics, the go compiler itself will
// enforce that you can't mix up your ID types. For example, a function with
// the signature:
//
//	func f(id UserID) {}
//
// Will fail to compile if passed an id of type AccountID.
package typed
