package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	username = flag.String("u", os.Getenv("AUTH_USERNAME"), "username")
	password = flag.String("p", os.Getenv("AUTH_PASSWORD"), "password")
)

func hash(s string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return bytes, err
}

func main() {
	flag.Parse()
	if *username == "" {
		log.Fatal("не указано имя пользователя")
	}
	if *password == "" {
		log.Fatal("не указан пароль")
	}

	usernameHash, _ := hash(*username)
	passwordHash, _ := hash(*password)

	fmt.Println("имя пользователя: ", string(usernameHash))
	fmt.Println("пароль: ", string(passwordHash))
}
