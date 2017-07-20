package controller

import (
	"errors"
	"github.com/VickMellon/test-social-tournament/model"
	"github.com/VickMellon/test-social-tournament/storage"
)

var (
	ErrTournamentExists   error = errors.New("Tournament already exists")
	ErrTournamentFinished error = errors.New("Tournament already finished")
	ErrAlreadyJoined      error = errors.New("Player already joined to this tournament")
)

func (c *Controller) CreateTournament(tournamentId string, deposit float64) error {
	t := &model.Tournament{
		TournamentId: tournamentId,
		Deposit:      deposit,
		Status:       model.TOURNAMENT_STATUS_ANNOUNCED,
	}
	err := c.storage.SaveTournament(t)
	if err != nil {
		if err == storage.ErrAlreadyExists {
			return ErrTournamentExists
		} else {
			return err
		}
	}
	return nil
}

func (c *Controller) JoinTournament(tournamentId string, playerId string, backerIds []string) error {
	err := c.storage.JoinTournament(tournamentId, playerId, backerIds)
	if err != nil {
		if err == storage.ErrNotFound {
			return ErrPlayerNotFound
		} else if err == storage.ErrNotUpdated {
			return ErrLowBalance
		} else if err == storage.ErrAlreadyExists {
			return ErrAlreadyJoined
		} else {
			return err
		}
	}
	return nil
}

func (c *Controller) FinishTournament(tournamentId string, winners []*model.Winner) error {
	err := c.storage.FinishTournament(tournamentId, winners)
	if err != nil {
		if err == storage.ErrNotFound {
			return ErrPlayerNotFound
		} else if err == storage.ErrNotUpdated {
			return ErrPlayerNotFound
		} else if err == storage.ErrAlreadyExists {
			return ErrTournamentFinished
		} else {
			return err
		}
	}
	return nil
}
