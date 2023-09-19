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
	userID, _ := typeid.New[UserID]()
	accountID, _ := typeid.New[AccountID]()
	// Each ID should have the correct type prefix:
	fmt.Printf("User ID prefix: %s\n", userID.Type())
	fmt.Printf("Account ID prefix: %s\n", accountID.Type())
	// The compiler considers their go types to be different:
	fmt.Printf("%T != %T\n", userID, accountID)

	// Output:
	// User ID prefix: user
	// Account ID prefix: account
	// typed.TypeID[go.jetpack.io/typeid/typed_test.UserID] != typed.TypeID[go.jetpack.io/typeid/typed_test.AccountID]
}
