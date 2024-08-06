package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wilsonle/excalidraw-collab/models"
	"github.com/wilsonle/excalidraw-collab/pkg/database"
)

func main() {
	database.InitDB("./excalidraw-collab.db")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Do you want to create an admin user? (yes/no): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.TrimSpace(strings.ToLower(response))

	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	password = strings.TrimSpace(password)

	var role string
	if response == "yes" {
		role = "admin"
	} else {
		role = ""
	}

	err = models.CreateUser(username, password, role)
	if err != nil {
		log.Fatalf("Error creating user: %v\n", err)
	}

	if role == "admin" {
		fmt.Println("Admin user created successfully!")
	} else {
		fmt.Println("User created successfully without a role!")
	}
}
