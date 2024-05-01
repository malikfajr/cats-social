package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/malikfajr/cats-social/helper"
)

type Cat struct {
	Id          string    `json:"id"`
	UserEmail   string    `json:"userEmail"`
	Name        string    `json:"name"`
	Race        string    `json:"race"`
	Sex         string    `json:"sex"`
	AgeInMonth  int       `json:"ageInMonth"`
	ImageUrls   []string  `json:"imageUrls"`
	Description string    `json:"description"`
	HasMatched  bool      `json:"hasMatched"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CatInsertRequest struct {
	UserEmail   string   `json:"userEmail"`
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        string   `json:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         string   `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,dive,required,url"`
}

func SaveCat(ctx context.Context, tx *sql.Tx, cat CatInsertRequest) (int, time.Time) {
	SQL := "INSERT INTO cats (user_email, name, race, sex, age_in_month, image_urls, description) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at"
	id := 0
	var createdAt time.Time

	err := tx.QueryRowContext(ctx, SQL, cat.UserEmail, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, pq.Array(cat.ImageUrls), cat.Description).Scan(&id, &createdAt)
	helper.PanicIfError(err)

	createdAt.Format(time.RFC3339)

	return id, createdAt
}

type CatParam struct {
	Id            string
	Owned         string
	Email         string
	AgeStr        string
	HasMatchedStr string
	Race          string
	Sex           string
	Search        string
	Limit         string
	Offsset       string
}

func GetAllCat(ctx context.Context, tx *sql.Tx, catParam CatParam) []Cat {
	SQL := "SELECT id, name, race, sex, age_in_month, image_urls, description, hasmatched, created_at FROM cats WHERE TRUE"

	params := make([]interface{}, 0)

	if catParam.Owned != "" {
		owned, err := strconv.ParseBool(catParam.Owned)
		helper.PanicIfError(err)

		if owned == true {
			SQL += " AND user_email = $1"
		} else {
			SQL += " AND user_email IS != $1 "
		}

		params = append(params, catParam.Email)
	}

	if race := catParam.Race; race != "" {
		SQL += fmt.Sprintf(" AND race = $%d", len(params)+1)
		params = append(params, race)
	}

	if sex := catParam.Sex; sex != "" {
		SQL += fmt.Sprintf(" AND sex = $%d", len(params)+1)
		params = append(params, sex)
	}

	if hasMatched := catParam.HasMatchedStr; hasMatched != "" {
		match, err := strconv.ParseBool(hasMatched)
		helper.PanicIfError(err)

		SQL += fmt.Sprintf(" AND hasMatched = $%d", len(params)+1)
		params = append(params, match)
	}

	if ageStr := catParam.AgeStr; ageStr != "" {
		var operator string
		var ageValue int
		var err error

		switch ageStr[0] {
		case '>':
			operator = ">"
			ageValue, err = strconv.Atoi(ageStr[1:])
			helper.PanicIfError(err)
			break
		case '<':
			operator = "<"
			ageValue, err = strconv.Atoi(ageStr[1:])
			helper.PanicIfError(err)
			break
		default:
			operator = "="
			ageValue, err = strconv.Atoi(ageStr[1:])
			helper.PanicIfError(err)
		}

		SQL += fmt.Sprintf(" AND age_in_month %s $%d", operator, len(params)+1)
		params = append(params, ageValue)
	}

	if search := catParam.Search; search != "" {
		SQL += fmt.Sprintf(" AND name like $%d", len(params)+1)
		params = append(params, "%"+search+"%")
	}

	limit, err := strconv.Atoi(catParam.Limit)
	if err != nil {
		limit = 5
	}

	offset, err := strconv.Atoi(catParam.Offsset)
	if err != nil {
		offset = 0
	}

	SQL += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := tx.QueryContext(ctx, SQL, params...)
	helper.PanicIfError(err)
	defer rows.Close()

	cats := []Cat{}
	for rows.Next() {
		cat := &Cat{}
		rows.Scan(&cat.Id, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, pq.Array(&cat.ImageUrls), &cat.Description, &cat.HasMatched, &cat.CreatedAt)
		cat.CreatedAt.Format(time.RFC3339)

		cats = append(cats, *cat)
	}

	return cats
}

func GetCatById(ctx context.Context, tx *sql.Tx, Id int) (Cat, error) {
	cat := Cat{}
	SQL := "SELECT id, name, race, sex, age_in_month, image_urls, description, hasmatched, created_at FROM cats WHERE id = $1;"

	row, err := tx.QueryContext(ctx, SQL, Id)
	helper.PanicIfError(err)
	defer row.Close()

	if row.Next() == false {
		return cat, errors.New("cat id is not valid")
	}

	row.Scan(&cat.Id, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, &cat.ImageUrls, &cat.Description, &cat.HasMatched, &cat.CreatedAt)
	log.Println(cat)
	cat.CreatedAt.Format(time.RFC3339)
	return cat, nil
}

func DestroyCat(ctx context.Context, tx *sql.Tx, Id int) {
	SQL := "DELETE FROM cats WHERE id = $1;"

	_, err := tx.ExecContext(ctx, SQL, Id)
	helper.PanicIfError(err)
}