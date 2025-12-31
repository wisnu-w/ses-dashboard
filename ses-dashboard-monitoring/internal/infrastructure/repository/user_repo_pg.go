package repository

import (
	"context"
	"database/sql"
	"ses-monitoring/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `SELECT id, username, password, email, role, active FROM users WHERE username = $1 AND active = true`
	row := r.db.QueryRowContext(ctx, query, username)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Role, &u.Active)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*user.User, error) {
	query := `SELECT id, username, password, email, role, active FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Role, &u.Active)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (username, password, email, role, active) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, u.Username, u.Password, u.Email, u.Role, true)
	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*user.User, error) {
	query := `SELECT id, username, email, role, active FROM users ORDER BY username`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := &user.User{}
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.Active)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int, password string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, password, id)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id int, active bool) error {
	query := `UPDATE users SET active = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, active, id)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}