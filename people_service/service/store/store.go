package store

import (
	"context"
	"database/sql"

	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type Store struct {
	conn *pgx.Conn
}
type People struct {
	ID   int
	Name string
}

// NewStore creates new database connection
func NewStore(connString string) (*Store, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	return &Store{
		conn: conn,
	}, nil
}

// Не нашел как использовать pgx.Connect для функции WithInstance,
// так как первым параметром должен быть *sql.DB, а в пакете pgx нет такого поля
func runDbMigrations(conn string) error {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: "bragin",
		SchemaName:   "public",
	})
	if err != nil {
		return fmt.Errorf("migrate instance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "bragin", driver)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// Сначала получаю кол-во пользователей в таблице,
// чтобы не переалоцировать слайс, не знаю насколько это целеобразно
func (s *Store) ListPeople() ([]People, error) {
	var count int
	err := s.conn.QueryRow(context.Background(), `SELECT id FROM people
	ORDER BY id
	DESC LIMIT 1`).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("no entries in the table: %w", err)
	}
	people := make([]People, 0, count)
	rows, err := s.conn.Query(context.Background(), "SELECT id, name FROM people")
	if err != nil {
		return nil, fmt.Errorf("people query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		p := People{}
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, fmt.Errorf("scan failed: %v\n", err)
		}
		people = append(people, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("scan failed: %v\n", err)
	}
	return people, nil
}

func (s *Store) GetPeopleByID(id string) (People, error) {
	p := People{}
	err := s.conn.QueryRow(context.Background(), "SELECT * FROM people WHERE id="+id).Scan(&p.ID, &p.Name)
	if err != nil {
		return p, fmt.Errorf("people not found: %w", err)
	}
	return p, nil
}
