package pgstore

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"iryzzh/practicum-gophermart/internal/app/store"
)

type Store struct {
	db *sql.DB
	m  *migrate.Migrate
}

func New(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/pg",
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	s := &Store{
		db: db,
		m:  m,
	}

	if err := s.Up(); err != nil {
		return nil, err
	}

	return s, db.Ping()
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Up() error {
	if err := s.m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *Store) Down() error {
	if err := s.Close(); err != nil {
		return err
	}

	if err := s.m.Down(); err != nil {
		return err
	}

	return nil
}

func (s *Store) User() store.UserRepository {
	return &UserRepository{store: s}
}

func (s *Store) Order() store.OrderRepository {
	return &OrderRepository{store: s}
}

func (s *Store) Balance() store.BalanceRepository {
	return &BalanceRepository{store: s}
}
