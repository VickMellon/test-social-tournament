package controller

import "github.com/VickMellon/test-social-tournament/storage"

type Controller struct {
	storage *storage.Storage
}

func NewController(storage *storage.Storage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (c *Controller) ResetDB() error {
	return c.storage.Reset()
}
