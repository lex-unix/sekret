package sekret_test

import (
	"fmt"

	"github.com/lex-unix/sekret"
)

func ExampleNew() {
	password := sekret.New("supersecret123")

	// Safe printing - secret is masked
	fmt.Printf("Password: %s\n", password)

	// Deliberate access when you need the actual value
	fmt.Printf("Actual value: %s\n", password.ExposeSecret())

	// Output:
	// Password: ******
	// Actual value: supersecret123
}
