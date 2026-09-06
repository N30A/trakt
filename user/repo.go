package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createUserLockID int64 = 1

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
	ErrInternal = errors.New("internal error")
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool}
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (User, error) {
	query := `
		SELECT
			id, email, password_hash, role
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT
			id, email, password_hash, role
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, newUser User) (User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", createUserLockID)
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users)").Scan(&exists) // always returns true or false, no need to check for ErrNoRows
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role
	`

	role := newUser.Role
	if !exists {
		role = RoleAdmin
	}

	var user User
	row := tx.QueryRow(ctx, query, newUser.Email, newUser.PasswordHash, role)
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return User{}, fmt.Errorf("%w: %v", ErrConflict, err)
		}
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}
