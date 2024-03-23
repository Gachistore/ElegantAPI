package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, email, pw string, userType UserType) *Account {
	acc, err := NewAccount(fname, lname, email, pw, userType)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("new account ", acc)
	return acc
}

func seedAccounts(store Storage) {
	seedAccount(store, "a", "b", "a@b.com", "1337", UserTypeAdmin)
	seedAccount(store, "c", "d", "c@d.com", "1337", UserTypeRegular)
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	for _, err1 := range store.Init() {
		if err1 != nil {
			log.Fatal(err1)
		}
	}
	if *seed {
		fmt.Println("seeding the db")
		seedAccounts(store)
	}

	server := newAPIServer(":3000", store)
	server.Run()
}
