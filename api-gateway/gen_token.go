package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	token := jwt.New(jwt.SigningMethodHS256)
	s, _ := token.SignedString([]byte("securepay-secret-key"))
	fmt.Println(s)
}
