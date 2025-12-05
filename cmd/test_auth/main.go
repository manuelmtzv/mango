package main

import (
	"fmt"
	"log"

	"github.com/manuelmtzv/mangocatnotes-api/internal/auth"
)

func main() {
	password := "password123"
	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	fmt.Printf("Hash: %s\n", hash)

	match, err := auth.VerifyPassword(password, hash)
	if err != nil {
		log.Fatalf("Failed to verify password: %v", err)
	}

	if match {
		fmt.Println("Password match!")
	} else {
		fmt.Println("Password mismatch!")
	}

	// Test with wrong password
	match, err = auth.VerifyPassword("wrongpassword", hash)
	if err != nil {
		log.Fatalf("Failed to verify wrong password: %v", err)
	}

	if !match {
		fmt.Println("Wrong password correctly rejected")
	} else {
		fmt.Println("Wrong password incorrectly accepted!")
	}
}
