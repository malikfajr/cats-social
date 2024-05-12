package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/malikfajr/cats-social/helper"
)

type ApproveRemoveRequest struct {
	MatchId string `json:"matchId" validate:"required"`
}

type Issuer struct {
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"createdAt"`
}

func (d Issuer) toJson() ([]byte, error) {
	b, err := json.Marshal(d)
	return b, err
}

type CatDetail struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Race        string     `json:"race"`
	Sex         string     `json:"sex"`
	Description string     `json:"description"`
	AgeInMonth  int        `json:"ageInMonth"`
	ImageUrls   []string   `json:"imageUrls"`
	HasMatched  bool       `json:"hasMatched"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func (d CatDetail) toJson() ([]byte, error) {
	b, err := json.Marshal(d)
	return b, err
}

type Match struct {
	Id             string    `json:"id"`
	IssuedBy       Issuer    `json:"issuedBy"`
	MatchCatDetail CatDetail `json:"matchCatDetail"`
	MatchUserId    int       `json:"-"`
	Status         string    `json:"-"`
	UserCatDetail  CatDetail `json:"userCatDetail"`
	Message        string    `json:"message"`
	CreatedAt      string    `json:"createdAt"`
}

type MatchInsertPayload struct {
	IssuedUserId int
	MatchUserId  int
	UserCatId    string
	MatchCatId   string
	Message      string
}

type MatchInsertRequest struct {
	MatchCatId string `json:"matchCatId" validate:"required"`
	UserCatId  string `json:"userCatId" validate:"required"`
	Message    string `json:"message" validate:"required,min=5,max=120"`
}

func NewMatch(ctx context.Context, tx *sql.Tx, match *MatchInsertPayload) (string, *time.Time, error) {
	var id string = ""
	var createdAt *time.Time
	SQL := "INSERT INTO matches (status, issued_user_id, match_user_id, user_cat_id, match_cat_id, message) VALUES ('pending', $1, $2, $3, $4, $5) RETURNING id, created_at"
	err := tx.QueryRowContext(ctx, SQL, match.IssuedUserId, match.MatchUserId, match.UserCatId, match.MatchCatId, match.Message).Scan(&id, &createdAt)

	helper.PanicIfError(err)

	return id, createdAt, err
}

func CrossCheckMatchCatId(ctx context.Context, tx *sql.Tx, matchCatId string, userCatId string) int {
	count := 0
	SQL := "SELECT 1 FROM matches WHERE (match_cat_id = $1 AND user_cat_id = $2) OR (match_cat_id = $3 AND user_cat_id = $4) LIMIT 1;"

	err := tx.QueryRowContext(ctx, SQL, matchCatId, userCatId, userCatId, matchCatId).Scan(&count)
	if err != nil {
		return 0
	}

	return 1
}

func GetAllMatch(ctx context.Context, db *sql.DB, userId int) ([]Match, error) {
	matches := []Match{}
	SQL := `
		SELECT m.id, 
			(
				SELECT json_build_object(
					'name', issuer.name, 
					'email', issuer.email, 
					'createdAt', to_char(issuer.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'))
				FROM users AS issuer WHERE issuer.id = m.issued_user_id)  AS issuedBy ,
			(
				SELECT json_build_object(
					'id', mc.id,
					'name', mc.name,
					'race', mc.race,
					'sex', mc.sex,
					'description', mc.description,
					'ageInMonth', mc.age_in_month,
					'imageUrls', mc.image_urls,
					'hasMatched', mc.hasMatched,
					'createdAt', to_char(mc.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
				)
				FROM cats AS mc WHERE mc.id = m.match_cat_id
			) AS matchCatDetail,
				(
				SELECT json_build_object(
					'id', uc.id,
					'name', uc.name,
					'race', uc.race,
					'sex', uc.sex,
					'description', uc.description,
					'ageInMonth', uc.age_in_month,
					'imageUrls', uc.image_urls,
					'hasMatched', uc.hasMatched,
					'createdAt', to_char(uc.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
				)
				FROM cats AS uc WHERE uc.id = m.match_cat_id
			) AS userCatDetail,
			m.message,
			m.created_at
		FROM matches m
		WHERE m.status = 'pending' AND m.issued_user_id = $1 OR m.match_user_id = $2
	`

	rows, err := db.QueryContext(ctx, SQL, userId, userId)
	if err != nil {
		return matches, err
	}
	defer rows.Close()

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

func DeleteMatch(ctx context.Context, tx *sql.Tx, matchId string) (string, int, error) {
	var idStr, status string
	var issuedUserId int
	SQL := "DELETE FROM matches WHERE id = $1 RETURNING id, status, issued_user_id"

	err := tx.QueryRowContext(ctx, SQL, matchId).Scan(&idStr, &status, &issuedUserId)

	return status, issuedUserId, err
}

func ApproveMatch(ctx context.Context, tx *sql.Tx, matchId string) {
	SQL := `UPDATE matches 
				SET status = 'approved', 
					match_cat_detail=jsonb_set(match_cat_detail, '{hasMatched}', 'true'), 
					user_cat_detail=jsonb_set(user_cat_detail, '{hasMatched}', 'true') 
			WHERE id = $1`

	tx.QueryRowContext(ctx, SQL, matchId).Scan()
}

func RejectMatch(ctx context.Context, tx *sql.Tx, matchId string) {
	SQL := "UPDATE matches SET status = 'reject' WHERE id = $1"

	tx.QueryRowContext(ctx, SQL, matchId).Scan()
}

func RejectOtherMatch(ctx context.Context, tx *sql.Tx, catId1 string, catId2 string, matchId string) {
	SQL1 := `UPDATE matches 
				SET status = 'reject'
			WHERE id != $1 
			AND status != 'approved' 
			AND user_cat_id = $2
			OR user_cat_id = $3
			OR match_cat_id = $4
			OR match_cat_id = $5
			;`

	_, err := tx.ExecContext(ctx, SQL1, matchId, catId1, catId2, catId1, catId2)
	helper.PanicIfError(err)
}

func GetMatchById(ctx context.Context, tx *sql.Tx, matchId string) (Match, error) {
	match := &Match{}
	var issuedByStr, matchCatDetailStr, userCatDetailStr string
	SQL := `
	SELECT m.id, 
			m.match_user_id,
			m.status,
			(
				SELECT json_build_object(
					'name', issuer.name, 
					'email', issuer.email, 
					'createdAt', to_char(issuer.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'))
				FROM users AS issuer WHERE issuer.id = m.issued_user_id)  AS issuedBy ,
			(
				SELECT json_build_object(
					'id', mc.id,
					'name', mc.name,
					'race', mc.race,
					'sex', mc.sex,
					'description', mc.description,
					'ageInMonth', mc.age_in_month,
					'imageUrls', mc.image_urls,
					'hasMatched', mc.hasMatched,
					'createdAt', to_char(mc.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
				)
				FROM cats AS mc WHERE mc.id = m.match_cat_id
			) AS matchCatDetail,
				(
				SELECT json_build_object(
					'id', uc.id,
					'name', uc.name,
					'race', uc.race,
					'sex', uc.sex,
					'description', uc.description,
					'ageInMonth', uc.age_in_month,
					'imageUrls', uc.image_urls,
					'hasMatched', uc.hasMatched,
					'createdAt', to_char(uc.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
				)
				FROM cats AS uc WHERE uc.id = m.match_cat_id
			) AS userCatDetail,
			m.message,
			m.created_at
		FROM matches m
		WHERE m.id = $1 LIMIT 1
	`

	err := tx.QueryRowContext(ctx, SQL, matchId).Scan(&match.Id, &match.MatchUserId, &match.Status, &issuedByStr, &matchCatDetailStr, &userCatDetailStr, &match.Message, &match.CreatedAt)

	if err == nil {
		// text to json/struct
		err := json.Unmarshal([]byte(issuedByStr), &match.IssuedBy)
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
	SQL := "SELECT 1 FROM matches WHERE user_cat_id = $1 OR match_cat_id = $2"

	tx.QueryRowContext(ctx, SQL, catId, catId).Scan(&count)

	return count
}
