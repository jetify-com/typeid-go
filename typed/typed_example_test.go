package typed_test

import (
	"fmt"

	typeid "go.jetpack.io/typeid/typed"
)

type userPrefix struct{}

func (userPrefix) Type() string { return "user" }

type UserID struct{ typeid.TypeID[userPrefix] }

type accountPrefix struct{}

func (accountPrefix) Type() string { return "account" }

type AccountID struct{ typeid.TypeID[accountPrefix] }

func Example() {
	user_id, _ := typeid.New[UserID]()
	account_id, _ := typeid.New[AccountID]()
	// Each ID should have the correct type prefix:
	fmt.Printf("User ID prefix: %s\n", user_id.Type())
	fmt.Printf("Account ID prefix: %s\n", account_id.Type())
	// The compiler considers their go types to be different:
	fmt.Printf("%T != %T\n", user_id, account_id)

	// Output:
	// User ID prefix: user
	// Account ID prefix: account
	// typed.TypeID[go.jetpack.io/typeid/typed_test.UserID] != typed.TypeID[go.jetpack.io/typeid/typed_test.AccountID]
}
