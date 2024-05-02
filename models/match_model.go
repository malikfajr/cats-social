package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/malikfajr/cats-social/helper"
)

type Issuer struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func (d Issuer) toJson() ([]byte, error) {
	b, err := json.Marshal(d)
	return b, err
}

type CatDetail struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Race        string    `json:"race"`
	Sex         string    `json:"sex"`
	Description string    `json:"description"`
	AgeInMonth  int       `json:"ageInMonth"`
	ImageUrls   []string  `json:"imageUrls"`
	HasMatched  bool      `json:"hasMatched"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (d CatDetail) toJson() ([]byte, error) {
	b, err := json.Marshal(d)
	return b, err
}

type Match struct {
	Id             string    `json:"id"`
	IssuedBy       Issuer    `json:"issuedBy"`
	MatchCatDetail CatDetail `json:"matchCatDetail"`
	MatchUserEmail string    `json:"-"`
	Status         string    `json:"-"`
	UserCatDetail  CatDetail `json:"userCatDetail"`
	Message        string    `json:"string"`
	CreatedAt      time.Time `json:"createdAt"`
}

type MatchInsertRequest struct {
	MatchCatId string `json:"matchCatId" validate:"required"`
	UserCatId  string `json:"userCatId" validate:"required"`
	Message    string `json:"message" validate:"required,min=5,max=120"`
}

func NewMatch(ctx context.Context, tx *sql.Tx, match Match) (string, error) {
	var id string = ""
	SQL := "INSERT INTO matches (match_user_email, issued_by, match_cat_detail, user_cat_detail, message) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	issuedBy, _ := match.IssuedBy.toJson()
	matchCat, _ := match.MatchCatDetail.toJson()
	userCat, _ := match.UserCatDetail.toJson()

	// err := tx.QueryRowContext(ctx, SQL, match.IssuedBy, pq.Array(match.MatchCatDetail), pq.Array(match.UserCatDetail), match.Message).Scan(&id)
	err := tx.QueryRowContext(ctx, SQL, match.MatchUserEmail, string(issuedBy), string(matchCat), string(userCat), match.Message).Scan(&id)
	helper.PanicIfError(err)

	return id, err
}

func CrossCheckMatchCatId(ctx context.Context, tx *sql.Tx, matchCatId string, userCatId string) int {
	count := 0
	SQL := "SELECT COUNT(*) WHERE match_cat_detail->>'id' = $1 AND user_cat_detail->>'id' = $2"

	tx.QueryRowContext(ctx, SQL, &count)

	return count
}

func GetAllMatch(ctx context.Context, tx *sql.Tx, email string) ([]Match, error) {
	matches := []Match{}
	SQL := "SELECT id, issued_by, match_cat_detail, user_cat_detail, message, created_at FROM matches WHERE issued_by->>'email' = $1 OR match_user_email = $2"

	rows, err := tx.QueryContext(ctx, SQL, email, email)
	if err != nil {
		return matches, err
	}

	for rows.Next() {
		match := &Match{}
		var issuedByStr, matchCatDetailStr, userCatDetailStr string

		err := rows.Scan(&match.Id, &issuedByStr, &matchCatDetailStr, &userCatDetailStr, &match.Message, &match.CreatedAt)
		if err != nil {
			log.Println("Error scanning row: ", err)
		}

		// text to json/struct
		err = json.Unmarshal([]byte(issuedByStr), &match.IssuedBy)
		if err != nil {
			log.Println("Error unmarshalling IssuedBy JSON:", err)
			continue
		}

		err = json.Unmarshal([]byte(matchCatDetailStr), &match.MatchCatDetail)
		if err != nil {
			log.Println("Error unmarshalling MatchCatDetail JSON:", err)
			continue
		}

		err = json.Unmarshal([]byte(userCatDetailStr), &match.UserCatDetail)
		if err != nil {
			log.Println("Error unmarshalling UserCatDetail JSON:", err)
			continue
		}

		matches = append(matches, *match)
	}

	return matches, nil
}

func DeleteMatch(ctx context.Context, tx *sql.Tx, id string) (string, string, error) {
	var idStr, status, email string
	SQL := "DELETE FROM matches WHERE id = $1 RETURNING id, status, issued_by->>'email'"

	err := tx.QueryRowContext(ctx, SQL, id).Scan(&idStr, &status, &email)

	return status, email, err
}

func ApproveMatch(ctx context.Context, tx *sql.Tx, matchId string) {
	SQL := "UPDATE matches SET status = 'approved', match_cat_detail->>'hasMatched' = TRUE, user_cat_detail->>'hasMatched' = TRUE WHERE id = $1"

	tx.QueryRowContext(ctx, SQL, matchId).Scan()

}

func RejectMatch(ctx context.Context, tx *sql.Tx, matchId string) {
	SQL := "UPDATE matches SET status = 'reject', match_cat_detail->>'hasMatched' = TRUE, user_cat_detail->>'hasMatched' = TRUE WHERE id = $1"

	tx.QueryRowContext(ctx, SQL, matchId).Scan()
}

func GetMatchById(ctx context.Context, tx *sql.Tx, matchId string) (Match, error) {
	match := &Match{}
	var issuedByStr, matchCatDetailStr, userCatDetailStr string
	SQL := "SELECT id, issued_by, match_cat_detail, user_cat_detail, message, created_at FROM matches WHERE id = $1"

	err := tx.QueryRowContext(ctx, SQL, matchId).Scan(&match.Id, &issuedByStr, &matchCatDetailStr, &userCatDetailStr, &match.Message, &match.CreatedAt)

	if err == nil {
		// text to json/struct
		err = json.Unmarshal([]byte(issuedByStr), &match.IssuedBy)
		if err != nil {
			log.Println("Error unmarshalling IssuedBy JSON:", err)
		}

		err = json.Unmarshal([]byte(matchCatDetailStr), &match.MatchCatDetail)
		if err != nil {
			log.Println("Error unmarshalling MatchCatDetail JSON:", err)
		}

		err = json.Unmarshal([]byte(userCatDetailStr), &match.UserCatDetail)
		if err != nil {
			log.Println("Error unmarshalling UserCatDetail JSON:", err)
		}
	}

	return *match, err
}

func CountCatInMatch(ctx context.Context, tx *sql.Tx, catId string) int {
	var count int
	SQL := "SELECT COUNT (*) FROM matches WHERE match_cat_detail->>'id' = $1 OR user_cat_detail->>'id' = $2"

	tx.QueryRowContext(ctx, SQL, catId).Scan(&count)

	return count
}
