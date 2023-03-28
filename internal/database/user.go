package repository

import (
	"context"
	"fmt"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
	Email          string `json:"email" db:"email"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
}

func (r *Repository) Login(ctx context.Context, email, hashedPassword string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, name, surname, email from users where email = $1 AND hashed_password = $2`, email, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&u.Id, &u.Name, &u.Surname, &u.Email)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}

func (r *Repository) AddNewUser(ctx context.Context, name, surname, email, hashedPassword string) (err error) {
	_, err = r.pool.Exec(ctx, `insert into users (name, surname, email, hashed_password) values ($1, $2, $3 , $4)`, name, surname, email, hashedPassword)
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
