package keygen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//GenerateRandomKey creates a random key
func GenerateRandomKey(length int) string {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", k)
	// return k
}

//GenerateToken generates a token based on an input string
func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))

	return base64.StdEncoding.EncodeToString(hash)
}
