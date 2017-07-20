package storage

import (
	"database/sql"
	"github.com/VickMellon/test-social-tournament/model"
	"github.com/lib/pq"
)

func (s *Storage) SaveTournament(t *model.Tournament) error {
	res, err := s.stmtInsertTournament.Exec(t.TournamentId, t.Deposit, t.Status)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrAlreadyExists // tournament already exists
	}
	return nil
}

func (s *Storage) GetTournament(tournamentId string) (*model.Tournament, error) {
	t := new(model.Tournament)
	t.TournamentId = tournamentId
	if err := s.stmtGetTournament.QueryRow(tournamentId).Scan(&t.Deposit, &t.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound // tournament not exists
		} else {
			return nil, err // some other SQL error
		}
	}
	return t, nil
}

func (s *Storage) JoinTournament(tournamentId string, playerId string, backerIds []string) error {
	var err error
	var res sql.Result
	var rowsCount int64

	// Get tournament
	t, err := s.GetTournament(tournamentId)
	if err != nil {
		return err // not found?
	}
	if t.Status != model.TOURNAMENT_STATUS_ANNOUNCED {
		return ErrNotFound // can't join to finished tournament
	}
	// Check player
	_, err = s.GetPlayer(playerId)
	if err != nil {
		return err // not found?
	}
	// Check backers
	for _, backerId := range backerIds {
		_, err := s.GetPlayer(backerId)
		if err != nil {
			return err // not found?
		}
	}

	// calc part of deposit for all participants (player + backers)
	depositShare := t.Deposit / (1 + float64(len(backerIds)))

	// start transaction
	tx, err := s.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // for emergency interrupts

	// 1. Withdraw deposit part from player
	txStmtWithdrawDeposit, err := tx.Prepare(sqlUpdatePlayerBalance)
	if err != nil {
		tx.Rollback()
		return err
	}
	res, err = txStmtWithdrawDeposit.Exec(playerId, -depositShare)
	if err != nil {
		txStmtWithdrawDeposit.Close()
		tx.Rollback()
		return err
	}
	rowsCount, err = res.RowsAffected()
	if rowsCount != 1 {
		txStmtWithdrawDeposit.Close()
		tx.Rollback()
		return ErrNotUpdated
	}

	// 2. Withdraw deposit part from backers
	for _, backerId := range backerIds {
		res, err = txStmtWithdrawDeposit.Exec(backerId, -depositShare)
		if err != nil {
			txStmtWithdrawDeposit.Close()
			tx.Rollback()
			return err
		}
		rowsCount, err = res.RowsAffected()
		if rowsCount != 1 {
			txStmtWithdrawDeposit.Close()
			tx.Rollback()
			return ErrNotUpdated
		}
	}
	txStmtWithdrawDeposit.Close()

	// 3. Join player to tournament
	txStmtInsertParticipant, err := tx.Prepare(sqlInsertParticipant)
	if err != nil {
		tx.Rollback()
		return err
	}
	res, err = txStmtInsertParticipant.Exec(tournamentId, playerId, playerId)
	if err != nil {
		txStmtInsertParticipant.Close()
		tx.Rollback()
		return err
	}
	rowsCount, err = res.RowsAffected()
	if rowsCount != 1 {
		txStmtInsertParticipant.Close()
		tx.Rollback()
		return ErrAlreadyExists
	}

	// 4. Join backers to tournament
	for _, backerId := range backerIds {
		res, err = txStmtInsertParticipant.Exec(tournamentId, playerId, backerId)
		if err != nil {
			txStmtInsertParticipant.Close()
			tx.Rollback()
			return err
		}
		rowsCount, err = res.RowsAffected()
		if rowsCount != 1 {
			txStmtInsertParticipant.Close()
			tx.Rollback()
			return ErrAlreadyExists
		}
	}
	txStmtInsertParticipant.Close()

	// Commit operation
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FinishTournament(tournamentId string, winners []*model.Winner) error {
	var err error
	var res sql.Result
	var rows *sql.Rows
	var rowsCount int64

	// Get tournament
	t, err := s.GetTournament(tournamentId)
	if err != nil {
		return err // not found?
	}
	if t.Status == model.TOURNAMENT_STATUS_FINISHED {
		return ErrAlreadyExists // tournament finished
	}
	// Check winners
	winnersIds := make([]string, 0)
	winnersMap := make(map[string]*model.Winner, 0)
	for _, w := range winners {
		winnersIds = append(winnersIds, w.PlayerId)
		winnersMap[w.PlayerId] = w
	}
	rows, err = s.stmtSelectTournamentWinners.Query(tournamentId, pq.Array(winnersIds))
	if err != nil {
		return err
	}
	// fetch results
	for rows.Next() {
		var playerId, backerId string
		err := rows.Scan(&playerId, &backerId)
		if err != nil {
			return err
		}
		if winner, exists := winnersMap[playerId]; exists && winner != nil {
			winner.IsParticipant = true // confirm participation
			// collect backers
			winner.Backers = append(winner.Backers, backerId)
		}
	}

	// start transaction
	tx, err := s.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // for emergency interrupts

	// 1. Finish tournament
	txStmtUpdateTournamentStatus, err := tx.Prepare(sqlUpdateTournamentStatus)
	if err != nil {
		tx.Rollback()
		return err
	}
	res, err = txStmtUpdateTournamentStatus.Exec(tournamentId, model.TOURNAMENT_STATUS_ANNOUNCED, model.TOURNAMENT_STATUS_FINISHED)
	if err != nil {
		txStmtUpdateTournamentStatus.Close()
		tx.Rollback()
		return err
	}
	rowsCount, err = res.RowsAffected()
	if rowsCount != 1 {
		txStmtUpdateTournamentStatus.Close()
		tx.Rollback()
		return ErrAlreadyExists
	}
	txStmtUpdateTournamentStatus.Close()

	// 2. Fund prize to win players and backers
	txStmtUpdatePlayerBalance, err := tx.Prepare(sqlUpdatePlayerBalance)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, winner := range winnersMap {
		if !winner.IsParticipant {
			continue // What, winner, but not participant? skip
		}
		backerIds := winner.Backers
		prizePart := winner.Prize / (float64(len(backerIds)))
		// reward player & backers
		for _, backerId := range backerIds {
			res, err = txStmtUpdatePlayerBalance.Exec(backerId, prizePart)
			if err != nil {
				txStmtUpdatePlayerBalance.Close()
				tx.Rollback()
				return err
			}
			rowsCount, err = res.RowsAffected()
			if rowsCount != 1 {
				txStmtUpdatePlayerBalance.Close()
				tx.Rollback()
				return ErrNotUpdated
			}
		}
	}
	txStmtUpdatePlayerBalance.Close()

	// Commit operation
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
