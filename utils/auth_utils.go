package utils

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Verify(user string, pass string) bool {
	// TODO: manage user system with database
	savePassword, readPasswordErr := os.ReadFile(user + ".txt")
	if readPasswordErr != nil {
		log.Fatal(readPasswordErr)
	}

	authPasswordStr := []byte(string(savePassword))
	authPasswordByte, _ := bcrypt.GenerateFromPassword(authPasswordStr, bcrypt.DefaultCost)

	hashCompareErr := bcrypt.CompareHashAndPassword(authPasswordByte, []byte(pass))
	if hashCompareErr != nil {
		fmt.Println(hashCompareErr)
		return false
	}
	return true
}
