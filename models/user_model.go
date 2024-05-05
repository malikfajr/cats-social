package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/malikfajr/cats-social/helper"
)

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

func SaveUser(ctx context.Context, tx *sql.Tx, user User) User {
	SQL := "INSERT INTO users (email, name, password) VALUES ($1, $2, $3);"

	_, err := tx.ExecContext(ctx, SQL, user.Email, user.Name, user.Password)
	helper.PanicIfError(err)

	return user
}

type Credential struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,string,min=5,max=15"`
}

func GetUserByEmail(ctx context.Context, tx *sql.Tx, email string) (User, error) {
	user := User{}

	SQL := "SELECT email, password, name FROM users WHERE email = $1;"

	row, err := tx.QueryContext(ctx, SQL, email)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	if row.Next() == false {
		return user, errors.New("user not found")
	}

	row.Scan(&user.Email, &user.Password, &user.Name)
	return user, nil
}
