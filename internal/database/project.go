package repository

import (
	"context"
	"fmt"
)

type Project struct {
	Id          int    `json:"id" db:"projectid"`
	CreatedAt   string `json:"createdat" db:"createdat"`
	ProjectName string `json:"projectname" db:"projectname"`
	Comment     string `json:"comment" db:"comment"`
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

func (r *Repository) CreateProject(ctx context.Context, projectname, email, comment string) (err error) {
	_, err = r.pool.Exec(ctx, `insert into projects (projectname, author, comment, author) 
	values ($1, $2, $3, (select id from users where email = $2) )`, projectname, email, comment)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}

	return
}
func (r *Repository) DeleteProject(ctx context.Context, projectname string) (err error) {
	_, err = r.pool.Exec(ctx, `delete from projects where projectname=$1)`, projectname)
	if err != nil {
		err = fmt.Errorf("failed to exec data: %w", err)
		return
	}

	return
}
