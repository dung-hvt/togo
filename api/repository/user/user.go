package user

import (
	"context"
	"database/sql"

	"manabie/todo/models"
)

const (
	queryFind = `SELECT * FROM member`
)

type UserRespository interface {
	Find(ctx context.Context, tx *sql.Tx) ([]*models.User, error)
	Create(ctx context.Context, tx *sql.Tx, u *models.User) error
}

type userRespository struct{}

func NewUserRespository() UserRespository {
	return &userRespository{}
}

func (ur *userRespository) Find(ctx context.Context, tx *sql.Tx) ([]*models.User, error) {
	rows, err := tx.QueryContext(ctx, queryFind)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	users := make([]*models.User, 0)

	for rows.Next() {

		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdateAt); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (ur *userRespository) Create(ctx context.Context, tx *sql.Tx, u *models.User) error {
	var (
		query string
		args  []interface{}
	)

	if u.ID != 0 {
		query, args = `INSERT INTO member (id, email, name) VALUES ($1, $2, $3)`, []interface{}{u.ID, u.Email, u.Name}
	} else {
		query, args = `INSERT INTO member (email, name) VALUES ($1, $2)`, []interface{}{u.Email, u.Name}
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return err
	}

	return nil
}
