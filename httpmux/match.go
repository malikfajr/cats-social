package httpmux

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/malikfajr/cats-social/exception"
	"github.com/malikfajr/cats-social/helper"
	"github.com/malikfajr/cats-social/models"
)

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)

	matchBody := models.MatchInsertRequest{}

	json.NewDecoder(r.Body).Decode(&matchBody)

	err := validate.Struct(matchBody)
	helper.PanicIfError(err)

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	issuerCat, err := models.GetCatById(r.Context(), tx, matchBody.UserCatId)
	if err != nil {
		panic(exception.NewBadRequestError("user cat id not found"))
	}

	if issuerCat.UserId != userId {
		panic(exception.NewNotFoundError("userCatId is not belong to the user"))
	}

	receiverCat, err := models.GetCatById(r.Context(), tx, matchBody.MatchCatId)
	if err != nil {
		panic(exception.NewBadRequestError("match cat id not found"))
	}

	if issuerCat.Sex == receiverCat.Sex {
		panic(exception.NewBadRequestError("gender cannot same"))
	}

	if issuerCat.UserId == receiverCat.UserId {
		panic(exception.NewBadRequestError("cannot match the same owner"))
	}

	matchInsert := &models.MatchInsertPayload{
		IssuedUserId: userId,
		MatchUserId:  receiverCat.UserId,
		UserCatId:    matchBody.UserCatId,
		MatchCatId:   matchBody.MatchCatId,
		Message:      matchBody.Message,
	}

	exist := models.CrossCheckMatchCatId(r.Context(), tx, matchInsert.MatchCatId, matchBody.UserCatId)
	if exist != 0 {
		panic(exception.NewBadRequestError("Cat id already submit to match"))
	}

	id, createdAt, err := models.NewMatch(r.Context(), tx, matchInsert)

	wrapper := &helper.WebResponse{
		Message: "success",
		Data: map[string]string{
			"matchId":   id,
			"createdAt": createdAt.Format(time.RFC3339),
		},
	}

	helper.WriteToResponseBody(w, *wrapper, http.StatusCreated)
}

func GetMyMatch(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)

	db := models.GetDb()

	matches, err := models.GetAllMatch(r.Context(), db, userId)
	helper.PanicIfError(err)

	wrapper := &helper.WebResponse{
		Message: "success",
		Data:    matches,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusOK)
}

func ApproveMatch(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)
	var bodyRequest models.ApproveRemoveRequest

	json.NewDecoder(r.Body).Decode(&bodyRequest)

	err := validate.Struct(bodyRequest)
	helper.PanicIfError(err)

	matchId := bodyRequest.MatchId

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	match, err := models.GetMatchById(r.Context(), tx, matchId)
	if err != nil {
		panic(exception.NewNotFoundError("match id not found"))
	}

	// check valid userId, if the user is not valid receiver match
	if userId != match.MatchUserId {
		panic(exception.NewNotFoundError("match id not found"))
	}

	if match.Status != "pending" {
		panic(exception.NewBadRequestError("match id is no longer valid"))
	}

	models.ApproveMatch(r.Context(), tx, matchId)
	models.RejectOtherMatch(r.Context(), tx, match.MatchCatDetail.Id, match.UserCatDetail.Id, matchId)

	models.UpdateStatusCat(r.Context(), tx, match.MatchCatDetail.Id, match.UserCatDetail.Id)

	helper.WriteToResponseBody(w, nil, http.StatusOK)
}

func RejectMatch(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)
	matchId := r.PathValue("id")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	match, err := models.GetMatchById(r.Context(), tx, matchId)
	if err != nil {
		panic(exception.NewNotFoundError("match id not found"))
	}

	// check valid email, if the user is not valid receiver match
	if userId != match.MatchUserId {
		panic(exception.NewNotFoundError("match id not found"))
	}

	if match.Status != "pending" {
		panic(exception.NewBadRequestError("match id is no longer valid"))
	}

	models.RejectMatch(r.Context(), tx, matchId)

	helper.WriteToResponseBody(w, nil, http.StatusOK)
}

func DeleteMatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId := r.Context().Value("user.id").(int)

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	status, issuedUserId, err := models.DeleteMatch(r.Context(), tx, id)
	if err != nil {
		panic(exception.NewNotFoundError("match id not found"))
	}

	if issuedUserId != userId {
		panic(exception.NewBadRequestError("You not issuer"))
	}

	if status != "pending" {
		panic(exception.NewBadRequestError("match is already approved / reject"))
	}

	helper.WriteToResponseBody(w, nil, http.StatusOK)
}
