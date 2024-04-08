package main

import (
	"3legant/api"
	"3legant/storage"
	"3legant/types"
	"flag"
	"fmt"
	"log"
)

func seedAccount(store storage.Storage, fname, lname, email, pw string, userType types.UserType) *types.Account {
	acc, err := types.NewAccount(fname, lname, email, pw, userType)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("new account ", acc)
	return acc
}

func seedAccounts(store storage.Storage) {
	seedAccount(store, "a", "b", "a@b.com", "1337", types.UserTypeAdmin)
	seedAccount(store, "c", "d", "c@d.com", "1337", types.UserTypeRegular)
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()
	store, err := storage.NewPostgresStore()
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

	server := api.NewAPIServer(":3000", store)
	server.Run()
}
