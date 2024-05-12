package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/malikfajr/cats-social/exception"
)

type User struct {
	Id       int    `json:"-"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

func SaveUser(ctx context.Context, tx *sql.Tx, user User) User {
	SQL := "INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id;"

	err := tx.QueryRowContext(ctx, SQL, user.Email, user.Name, user.Password).Scan(&user.Id)
	if err != nil {
		panic(exception.NewConflictError("Email is exists"))
	}

	return user
}

type Credential struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,string,min=5,max=15"`
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (User, error) {
	user := User{}

	SQL := "SELECT id, email, password, name FROM users WHERE email = $1 LIMIT 1;"

	row, err := db.QueryContext(ctx, SQL, email)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	if row.Next() == false {
		return user, errors.New("user not found")
	}

	row.Scan(&user.Id, &user.Email, &user.Password, &user.Name)
	return user, nil
}
