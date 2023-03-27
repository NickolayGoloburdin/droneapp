package repository

import (
	"context"
	"fmt"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Username       string `json:"username" db:"username"`
	Email          string `json:"email" db:"email"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
}

func (r *Repository) Login(ctx context.Context, email, hashedPassword string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, username, email from users where email = $1 AND hashed_password = $2`, email, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&u.Id, &u.Username, &u.Email)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}

func (r *Repository) AddNewUser(ctx context.Context, username, email, hashedPassword string) (err error) {
	_, err = r.pool.Exec(ctx, `insert into users (name, surname, login, hashed_password) values ($1, $2, $3)`, username, email, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}

	return
}
func (r *Repository) DeleteUser(ctx context.Context, email string) (err error) {
	_, err = r.pool.Exec(ctx, `delete from users where email=$1)`, email)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}

	return
}
