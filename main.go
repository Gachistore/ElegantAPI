package main

import "log"

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	for _, err1 := range store.Init() {
		if err1 != nil {
			log.Fatal(err1)
		}
	}

	server := newAPIServer(":3000", store)
	server.Run()
}
