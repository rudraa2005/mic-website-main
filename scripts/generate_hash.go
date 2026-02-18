package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_hash.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash:", string(hash))
	fmt.Println("\nSQL to update password:")
	fmt.Printf("UPDATE users SET password_hash = '%s', updated_at = NOW() WHERE email = 'YOUR_EMAIL';\n", string(hash))
}
