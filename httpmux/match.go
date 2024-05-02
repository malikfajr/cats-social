package httpmux

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/malikfajr/cats-social/exception"
	"github.com/malikfajr/cats-social/helper"
	"github.com/malikfajr/cats-social/models"
)

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	username := r.Header.Get("name")
	matchBody := models.MatchInsertRequest{}

	json.NewDecoder(r.Body).Decode(&matchBody)

	err := validate.Struct(matchBody)
	helper.PanicIfError(err)

	issuerCatId, err := strconv.Atoi(matchBody.UserCatId)
	helper.PanicIfError(err)

	receiverCatId, err := strconv.Atoi(matchBody.MatchCatId)
	helper.PanicIfError(err)

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	issuerCat, err := models.GetCatById(r.Context(), tx, issuerCatId)
	if err != nil {
		panic(exception.NewNotFoundError("cat id not found"))
	}

	if issuerCat.UserEmail != email {
		panic(exception.NewNotFoundError("cat id not found"))
	}

	receiverCat, err := models.GetCatById(r.Context(), tx, receiverCatId)
	if err != nil {
		panic(exception.NewNotFoundError("cat id not found"))
	}

	if issuerCat.Sex == receiverCat.Sex {
		panic(exception.NewBadRequestError("gender cannot same"))
	}

	if issuerCat.UserEmail == receiverCat.UserEmail {
		panic(exception.NewBadRequestError("cannot match the same owner"))
	}

	matchInsert := models.Match{
		IssuedBy: models.Issuer{
			Email:     email,
			Name:      username,
			CreatedAt: time.Now(),
		},
		MatchUserEmail: receiverCat.UserEmail,
		MatchCatDetail: models.CatDetail{
			Id:          receiverCat.Id,
			Name:        receiverCat.Name,
			Race:        receiverCat.Race,
			Sex:         receiverCat.Sex,
			Description: receiverCat.Description,
			AgeInMonth:  receiverCat.AgeInMonth,
			ImageUrls:   receiverCat.ImageUrls,
			HasMatched:  receiverCat.HasMatched,
			CreatedAt:   receiverCat.CreatedAt,
		},
		UserCatDetail: models.CatDetail{
			Id:          issuerCat.Id,
			Name:        issuerCat.Name,
			Race:        issuerCat.Race,
			Sex:         issuerCat.Sex,
			Description: issuerCat.Description,
			AgeInMonth:  issuerCat.AgeInMonth,
			ImageUrls:   issuerCat.ImageUrls,
			HasMatched:  issuerCat.HasMatched,
			CreatedAt:   issuerCat.CreatedAt,
		},
		Message: matchBody.Message,
	}

	exist := models.CrossCheckMatchCatId(r.Context(), tx, issuerCat.Id, receiverCat.Id)
	if exist != 0 {
		panic(exception.NewNotFoundError("Cat id already submit to match"))
	}

	exist = models.CrossCheckMatchCatId(r.Context(), tx, receiverCat.Id, issuerCat.Id)
	if exist != 0 {
		panic(exception.NewNotFoundError("Cat id already submit to match"))
	}

	id, err := models.NewMatch(r.Context(), tx, matchInsert)

	wrapper := helper.WebResponse{
		Message: "success",
		Data: map[string]string{
			"matchId":   id,
			"createdAt": matchInsert.IssuedBy.CreatedAt.Format(time.RFC3339),
		},
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusCreated)
}

func GetMyMatch(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	matches, err := models.GetAllMatch(r.Context(), tx, email)
	helper.PanicIfError(err)

	wrapper := &helper.WebResponse{
		Message: "success",
		Data:    matches,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusOK)
}

func ApproveMatch(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	matchId := r.PathValue("id")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	match, err := models.GetMatchById(r.Context(), tx, matchId)
	if err != nil {
		panic(exception.NewNotFoundError("match id not found"))
	}

	// check valid email, if the user is not valid receiver match
	if email != match.MatchUserEmail {
		panic(exception.NewNotFoundError("match id not found"))
	}

	if match.Status != "pending" {
		panic(exception.NewBadRequestError("match id is no longer valid"))
	}

	models.ApproveMatch(r.Context(), tx, matchId)

	// TODO: Update cat property and other match status

	helper.WriteToResponseBody(w, nil, http.StatusOK)
}

func RejectMatch(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	matchId := r.PathValue("id")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	match, err := models.GetMatchById(r.Context(), tx, matchId)
	if err != nil {
		panic(exception.NewNotFoundError("match id not found"))
	}

	// check valid email, if the user is not valid receiver match
	if email != match.MatchUserEmail {
		panic(exception.NewNotFoundError("match id not found"))
	}

	if match.Status != "pending" {
		panic(exception.NewBadRequestError("match id is no longer valid"))
	}

	models.RejectMatch(r.Context(), tx, matchId)

	// TODO: Update cat property and other match status

	helper.WriteToResponseBody(w, nil, http.StatusOK)
}

func DeleteMatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	issuerEmail := r.Header.Get("email")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	status, email, err := models.DeleteMatch(r.Context(), tx, id)
	if err != nil {
		panic(exception.NewNotFoundError("match id not found"))
	}

	if email != issuerEmail {
		panic(exception.NewBadRequestError("You not issuer"))
	}

	if status != "pending" {
		panic(exception.NewBadRequestError("match is already approved / reject"))
	}

	helper.WriteToResponseBody(w, nil, http.StatusOK)
}
