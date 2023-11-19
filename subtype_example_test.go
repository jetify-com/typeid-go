package typeid_test

import (
	"fmt"

	"go.jetpack.io/typeid"
)

// To create a new id type, first implement a custom PrefixType and ensure the
// Prefix() method returns the correct prefix
type UserPrefix struct{}

func (UserPrefix) Prefix() string {
	return "user"
}

// And then define your custom id type by embedding TypeID:
type UserID struct {
	typeid.TypeID[UserPrefix]
}

// Now do the same for AccountIDs
type AccountPrefix struct{}

func (AccountPrefix) Prefix() string {
	return "account"
}

type AccountID struct {
	typeid.TypeID[AccountPrefix]
}

func Example() {
	// To create new IDs call typeid.New and pass your custom id type as the
	// generic argument:
	userID, _ := typeid.New[UserID]()
	accountID, _ := typeid.New[AccountID]()

	// Other than that, your custom types should have the same methods as a
	// regular TypeID.
	// For example, we can check that each ID has the correct type prefix:
	fmt.Printf("User ID prefix: %s\n", userID.Prefix())
	fmt.Printf("Account ID prefix: %s\n", accountID.Prefix())

	// Despite both of them being TypeIDs, you now get compile-time safety because
	// the compiler considers their go types to be different:
	// (typeid_test.UserID vs typeid_test.AccountID vs typeid.TypeID)
	fmt.Printf("%T != %T\n", userID, accountID)
	// Output:
	// User ID prefix: user
	// Account ID prefix: account
	// typeid_test.UserID != typeid_test.AccountID
}
