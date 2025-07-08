package typeid_test

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"go.jetify.com/typeid"
)

// ExampleGenerate demonstrates creating a new TypeID for a user entity
func ExampleGenerate() {
	userID, err := typeid.Generate("user")
	if err != nil {
		panic(err)
	}

	// Use the generated ID for a new user
	fmt.Printf("Created new user with ID: %s\n", userID)
	// Output format: Created new user with ID: user_[26-char-suffix]
}

// ExampleParse demonstrates parsing a TypeID from a string, such as from a URL parameter
func ExampleParse() {
	// Parse a TypeID received from a client request
	orderIDStr := "order_00041061050r3gg28a1c60t3gf"
	orderID, err := typeid.Parse(orderIDStr)
	if err != nil {
		fmt.Printf("Invalid order ID: %v\n", err)
		return
	}

	fmt.Printf("Processing order: %s\n", orderID.String())
	fmt.Printf("Order type: %s\n", orderID.Prefix())
	// Output:
	// Processing order: order_00041061050r3gg28a1c60t3gf
	// Order type: order
}

// ExampleTypeID_MarshalText demonstrates using TypeID in JSON APIs
func ExampleTypeID_MarshalText() {
	// Define a struct for API responses
	type Product struct {
		ID    typeid.TypeID `json:"id"`
		Name  string        `json:"name"`
		Price float64       `json:"price"`
	}

	// Create a product with TypeID
	product := Product{
		ID:    typeid.MustParse("product_00041061050r3gg28a1c60t3gf"),
		Name:  "Widget",
		Price: 29.99,
	}

	// Marshal to JSON for API response
	jsonData, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jsonData)
	// Output:
	// {"id":"product_00041061050r3gg28a1c60t3gf","name":"Widget","price":29.99}
}

// ExampleTypeID_UnmarshalText demonstrates parsing TypeID from JSON requests
func ExampleTypeID_UnmarshalText() {
	// JSON payload from client request
	jsonPayload := `{
		"user_id": "user_00041061050r3gg28a1c60t3gf",
		"product_id": "product_00041061050r3gg28a1c60t3gg",
		"quantity": 2
	}`

	// Define request struct
	type OrderRequest struct {
		UserID    typeid.TypeID `json:"user_id"`
		ProductID typeid.TypeID `json:"product_id"`
		Quantity  int           `json:"quantity"`
	}

	// Parse the request
	var req OrderRequest
	err := json.Unmarshal([]byte(jsonPayload), &req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Order from user %s for product %s\n", req.UserID.Prefix(), req.ProductID.Prefix())
	// Output:
	// Order from user user for product product
}

// ExampleFromUUID demonstrates migrating from UUIDs to TypeIDs
func ExampleFromUUID() {
	// Existing UUID from legacy system
	existingUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Convert to TypeID with appropriate prefix
	customerID, err := typeid.FromUUID("customer", existingUUID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Migrated customer ID: %s\n", customerID.String())
	// Output:
	// Migrated customer ID: customer_2n1t201rmv87aae5j4csam8000
}

// ExampleTypeID_Scan demonstrates reading TypeIDs from a database
func ExampleTypeID_Scan() {
	// Simulate database row scan
	var userID typeid.TypeID

	// In real code, this would come from sql.Row.Scan()
	dbValue := "user_00041061050r3gg28a1c60t3gf"
	err := userID.Scan(dbValue)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Retrieved user %s from database\n", userID.String())
	// Output:
	// Retrieved user user_00041061050r3gg28a1c60t3gf from database
}

// Example_nullableColumns demonstrates using sql.Null[TypeID] for nullable database columns.
// This is the recommended approach for handling nullable TypeID columns in Go applications.
func Example_nullableColumns() {
	var managerID sql.Null[typeid.TypeID]

	// Scan NULL value from database
	err := managerID.Scan(nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Is valid: %v\n", managerID.Valid)

	// Scan actual TypeID value
	err = managerID.Scan("user_00041061050r3gg28a1c60t3gf")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Is valid: %v\n", managerID.Valid)
	fmt.Printf("Manager: %s\n", managerID.V.String())

	// Output:
	// Is valid: false
	// Is valid: true
	// Manager: user_00041061050r3gg28a1c60t3gf
}
