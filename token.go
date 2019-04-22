package token

import "fmt"
import "crypto/rand"

// Returns a string with random chars
func tokenGenerator() string {
	// Create the slice with
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
