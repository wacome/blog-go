package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "Twk20030415?"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
		return
	}
	fmt.Printf("Hashed password: %s\n", string(hashedPassword))
}
