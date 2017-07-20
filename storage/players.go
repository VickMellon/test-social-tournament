package storage

import (
	"database/sql"
	"github.com/VickMellon/test-social-tournament/model"
)

func (s *Storage) SavePlayer(player *model.Player) error {
	res, err := s.stmtInsertPlayer.Exec(player.PlayerId, player.Balance)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrAlreadyExists // player already exists
	}
	return nil
}

func (s *Storage) GetPlayer(playerId string) (*model.Player, error) {
	player := new(model.Player)
	player.PlayerId = playerId
	if err := s.stmtGetPlayerBalance.QueryRow(playerId).Scan(&player.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound // player not exists
		} else {
			return nil, err // some other SQL error
		}
	}
	return player, nil
}

func (s *Storage) UpdatePlayerBalance(playerId string, gain float64) error {
	res, err := s.stmtUpdatePlayerBalance.Exec(playerId, gain)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrNotUpdated // player not exists or balance + gain < 0
	}
	return nil
}
