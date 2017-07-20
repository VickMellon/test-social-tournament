package storage

import "fmt"

var (
	// Players
	sqlInsertPlayer = `
		INSERT INTO players
		(player_id, balance)
		VALUES
		($1, $2)
		ON CONFLICT DO NOTHING
	`
	sqlGetPlayerBalance = `
		SELECT balance
		FROM players
		WHERE player_id = $1
	`
	sqlUpdatePlayerBalance = `
		UPDATE players
		SET balance = balance + $2
		WHERE player_id = $1 AND balance + $2 >= 0
	`
	// Tournaments
	sqlInsertTournament = `
		INSERT INTO tournaments
		(tournament_id, deposit, status)
		VALUES
		($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	sqlGetTournament = `
		SELECT deposit, status
		FROM tournaments
		WHERE tournament_id = $1
	`
	sqlInsertParticipant = `
		INSERT INTO participants
		(tournament_id, player_id, backer_id)
		VALUES
		($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	sqlUpdateTournamentStatus = `
		UPDATE tournaments
		SET status = $3
		WHERE tournament_id = $1 AND status = $2
	`
	sqlSelectTournamentWinners = `
		SELECT player_id, backer_id
		FROM participants
		WHERE tournament_id = $1 AND player_id = ANY($2)
	`

	// Reset
	sqlDeleteAllParticipants = `
		TRUNCATE TABLE participants
	`
	sqlDeleteAllTournaments = `
		TRUNCATE TABLE tournaments
	`
	sqlDeleteAllPlayers = `
		TRUNCATE TABLE players
	`
)

func (s *Storage) initStatements() error {
	var err error

	// Players
	if s.stmtInsertPlayer, err = s.DB().Prepare(sqlInsertPlayer); err != nil {
		return fmt.Errorf("Prepare sqlInsertPlayer error: %v", err)
	}
	if s.stmtGetPlayerBalance, err = s.DB().Prepare(sqlGetPlayerBalance); err != nil {
		return fmt.Errorf("Prepare sqlGetPlayer error: %v", err)
	}
	if s.stmtUpdatePlayerBalance, err = s.DB().Prepare(sqlUpdatePlayerBalance); err != nil {
		return fmt.Errorf("Prepare sqlUpdatePlayerBalance error: %v", err)
	}

	// Tournaments
	if s.stmtInsertTournament, err = s.DB().Prepare(sqlInsertTournament); err != nil {
		return fmt.Errorf("Prepare sqlInsertTournament error: %v", err)
	}
	if s.stmtGetTournament, err = s.DB().Prepare(sqlGetTournament); err != nil {
		return fmt.Errorf("Prepare sqlGetTournament error: %v", err)
	}
	if s.stmtSelectTournamentWinners, err = s.DB().Prepare(sqlSelectTournamentWinners); err != nil {
		return fmt.Errorf("Prepare sqlSelectTournamentWinners error: %v", err)
	}

	// Reset
	if s.stmtDeleteAllParticipants, err = s.DB().Prepare(sqlDeleteAllParticipants); err != nil {
		return fmt.Errorf("Prepare sqlDeleteAllParticipants error: %v", err)
	}
	if s.stmtDeleteAllTournaments, err = s.DB().Prepare(sqlDeleteAllTournaments); err != nil {
		return fmt.Errorf("Prepare sqlDeleteAllTournaments error: %v", err)
	}
	if s.stmtDeleteAllPlayers, err = s.DB().Prepare(sqlDeleteAllPlayers); err != nil {
		return fmt.Errorf("Prepare sqlDeleteAllPlayers error: %v", err)
	}

	return nil
}
