package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	ErrNotFound      = errors.New("Record not found")
	ErrAlreadyExists = errors.New("Record already exists")
	ErrNotUpdated    = errors.New("No one record was updated")
)

type Storage struct {
	db *sql.DB

	stmtInsertPlayer        *sql.Stmt
	stmtGetPlayerBalance    *sql.Stmt
	stmtUpdatePlayerBalance *sql.Stmt

	stmtInsertTournament        *sql.Stmt
	stmtGetTournament           *sql.Stmt
	stmtSelectTournamentWinners *sql.Stmt

	stmtDeleteAllParticipants *sql.Stmt
	stmtDeleteAllTournaments  *sql.Stmt
	stmtDeleteAllPlayers      *sql.Stmt
}

type PostgresConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}

func NewStorage(conf *PostgresConfig) (*Storage, error) {
	var err error

	s := &Storage{}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", //&parseTime=True
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	if s.db, err = sql.Open("postgres", dsn); err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL DSN: %v", err)
	}

	if err = s.initStatements(); nil != err {
		s.Close()
		return nil, err
	}

	return s, nil
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) Close() error {
	db := s.db
	s.db = nil
	return db.Close()
}

func (s *Storage) Reset() error {
	//TODO Transaction needed?
	if _, err := s.stmtDeleteAllParticipants.Exec(); err != nil {
		return err
	}
	if _, err := s.stmtDeleteAllTournaments.Exec(); err != nil {
		return err
	}
	if _, err := s.stmtDeleteAllPlayers.Exec(); err != nil {
		return err
	}
	return nil
}
