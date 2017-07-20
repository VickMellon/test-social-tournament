package controller

import (
	"errors"
	"github.com/VickMellon/test-social-tournament/model"
	"github.com/VickMellon/test-social-tournament/storage"
)

var (
	ErrPlayerNotFound error = errors.New("Player not found")
	ErrLowBalance     error = errors.New("Player has not enough points")
)

func (c *Controller) TakePoints(playerId string, points float64) error {
	err := c.storage.UpdatePlayerBalance(playerId, -points)
	if err != nil {
		if err == storage.ErrNotUpdated {
			// player balance - points < 0
			return ErrLowBalance
		}
	}
	return nil
}

func (c *Controller) FundPoints(playerId string, points float64) error {
	err := c.storage.UpdatePlayerBalance(playerId, points)
	if err != nil {
		if err == storage.ErrNotUpdated {
			// player probably not exists
			// so try to create his
			return c.storage.SavePlayer(&model.Player{PlayerId: playerId, Balance: points})
		} else {
			return err
		}
	}
	return nil
}

func (c *Controller) GetBalance(playerId string) (float64, error) {
	player, err := c.storage.GetPlayer(playerId)
	if err != nil {
		if err == storage.ErrNotFound {
			return 0, ErrPlayerNotFound
		} else {
			return 0, err
		}
	}
	return player.Balance, nil
}
