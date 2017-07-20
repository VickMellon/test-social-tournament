package main

import (
	"github.com/VickMellon/test-social-tournament/controller"
	"github.com/VickMellon/test-social-tournament/server"
	"github.com/VickMellon/test-social-tournament/storage"
)

const (
	SERVICE_PORT  = ":8088"

	POSTGRES_HOST = "127.0.0.1"
	POSTGRES_PORT = 5432
	POSTGRES_USER = "test"
	POSTGRES_PASS = "test"
	POSTGRES_BASE = "test"
)

func main() {

	dbConf := &storage.PostgresConfig{
		Host:     POSTGRES_HOST,
		Port:     POSTGRES_PORT,
		User:     POSTGRES_USER,
		Password: POSTGRES_PASS,
		Database: POSTGRES_BASE,
	}
	s, err := storage.NewStorage(dbConf)
	if err != nil {
		panic(err)
	}
	c := controller.NewController(s)
	w := server.NewWebServer(c)

	w.Run(SERVICE_PORT)

	s.Close()
}
