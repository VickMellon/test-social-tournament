package main

import (
	"github.com/VickMellon/test-social-tournament/controller"
	"github.com/VickMellon/test-social-tournament/server"
	"github.com/VickMellon/test-social-tournament/storage"
)

func main() {

	dbConf := &storage.PostgresConfig{
		Host:     "192.168.0.99",
		Port:     5432,
		User:     "test",
		Password: "test",
		Database: "test",
	}
	s, err := storage.NewStorage(dbConf)
	if err != nil {
		panic(err)
	}
	c := controller.NewController(s)
	w := server.NewWebServer(c)

	w.Run(":8080")

	s.Close()
}