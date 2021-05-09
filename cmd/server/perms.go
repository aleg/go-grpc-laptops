package main

import (
	"log"

	"github.com/aleg/go-grpc-laptops/stores"
	"github.com/aleg/go-grpc-laptops/users"
)

func accessibleRoles() map[string][]string {
	path := "/aleg.laptops.LaptopService/"
	// SearchLaptop is accessible by everyone (even for
	// unregistered users).
	return map[string][]string{
		path + "CreateLaptop": {"admin"},
		path + "RateLaptop":   {"role1", "admin"},
		path + "UploadImage":  {}, // no user can access
	}

}

func createUsers(store *stores.InMemoryUserStore) {
	user1, _ := users.NewUser("jay", "secret-jay", "admin")
	user2, _ := users.NewUser("kay", "secret-kay", "role1")
	user3, _ := users.NewUser("rob", "secret-rob", "role2")

	log.Print("Creating user ", user1)
	log.Print("Creating user ", user2)
	log.Print("Creating user ", user3)

	store.Save(user1)
	store.Save(user2)
	store.Save(user3)
}
