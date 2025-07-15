package store

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/y3933y3933/joker/internal/db/sqlc"
	"github.com/y3933y3933/joker/internal/utils/errx"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}
	p.plainText = &plainText
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Password password `json:"-"`
}

type PostgresUserStore struct {
	queries *sqlc.Queries
}

func NewPostgresUserStore(queries *sqlc.Queries) *PostgresUserStore {
	return &PostgresUserStore{queries: queries}
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

func (pg *PostgresUserStore) Create(ctx context.Context, user *User) error {
	args := sqlc.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.Password.hash,
	}
	_, err := pg.queries.CreateUser(ctx, args)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return errx.ErrDuplicateUsername
			}
		}
		return err
	}
	return nil
}

func (pg *PostgresUserStore) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	row, err := pg.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	password := password{
		hash: row.PasswordHash,
	}

	return &User{
		ID:       row.ID,
		Username: username,
		Password: password,
	}, nil
}
