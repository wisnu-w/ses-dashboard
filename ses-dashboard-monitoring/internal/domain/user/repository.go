package user

import "context"

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, user *User) error
	GetAll(ctx context.Context) ([]*User, error)
	UpdatePassword(ctx context.Context, id int, password string) error
	UpdateStatus(ctx context.Context, id int, active bool) error
	Delete(ctx context.Context, id int) error
}