package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/guluzadehh/go_chat/internal/lib/db"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) UserByUsername(username string) (*models.User, error) {
	const op = "storage.sqlite.UserByUsername"

	var user models.User

	const query = `SELECT * FROM users WHERE username = ?`
	err := s.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.UserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) CreateUser(username, password string) (*models.User, error) {
	const op = "storage.sqlite.CreateUser"

	const query = `INSERT INTO users("username", "password") VALUES(?, ?)`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(username, password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, storage.UsernameExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.User{Id: lastInsertedId, Username: username, Password: password}, nil
}

func (s *Storage) UsersWithIds(ids []int64) (map[int64]*models.User, error) {
	const op = "storage.sqlite.UsersWithIds"

	query := fmt.Sprintf(`SELECT * FROM users WHERE users.id IN (%s)`, db.Placeholders(len(ids)))

	args := make([]interface{}, 0)
	for _, id := range ids {
		args = append(args, id)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	users := make(map[int64]*models.User)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users[user.Id] = &user
	}

	return users, nil
}
