package httpmux

import (
	"encoding/json"
	"net/http"

	"github.com/malikfajr/cats-social/exception"
	"github.com/malikfajr/cats-social/helper"
	"github.com/malikfajr/cats-social/models"
)

var RaceEnum = map[string]bool{
	"Persian":           true,
	"Maine Coon":        true,
	"Siamese":           true,
	"Ragdoll":           true,
	"Bengal":            true,
	"Sphynx":            true,
	"British Shorthair": true,
	"Abyssinian":        true,
	"Scottish Fold":     true,
	"Birman":            true,
}

var SexEnum = map[string]bool{
	"male":   true,
	"female": true,
}

func SaveCat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)
	catRequest := models.CatInsertRequest{}

	json.NewDecoder(r.Body).Decode(&catRequest)

	err := validate.Struct(catRequest)
	helper.PanicIfError(err)

	catRequest.UserId = userId

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)
	id, date := models.SaveCat(r.Context(), tx, catRequest)

	wraper := helper.WebResponse{
		Message: "success",
		Data: map[string]interface{}{
			"id":        id,
			"createdAt": date,
		},
	}

	helper.WriteToResponseBody(w, wraper, http.StatusCreated)
}

func GetCat(w http.ResponseWriter, r *http.Request) {
	var data []models.Cat = []models.Cat{}
	catParam := models.CatParam{
		Id:            r.URL.Query().Get("id"),
		Owned:         r.URL.Query().Get("owned"),
		UserId:        r.Context().Value("user.id").(int),
		AgeStr:        r.URL.Query().Get("ageInMonth"),
		HasMatchedStr: r.URL.Query().Get("hasMatched"),
		Race:          r.URL.Query().Get("race"),
		Sex:           r.URL.Query().Get("sex"),
		Search:        r.URL.Query().Get("search"),
		Limit:         r.URL.Query().Get("limit"),
		Offsset:       r.URL.Query().Get("offset"),
	}

	if ok := RaceEnum[catParam.Race]; ok == false {
		catParam.Race = ""
	}

	if ok := SexEnum[catParam.Sex]; ok == false {
		catParam.Sex = ""
	}

	db := models.GetDb()

	data = models.GetAllCat(r.Context(), db, catParam)

	wrapper := helper.WebResponse{
		Message: "success",
		Data:    data,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusOK)
}

func DestroyCat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)
	id := r.PathValue("id")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	err := models.DestroyCat(r.Context(), tx, id, userId)

	if err != nil {
		panic(exception.NewNotFoundError("id is not found"))
	}

	wrapper := helper.WebResponse{
		Message: "success",
		Data:    nil,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusOK)
}

func UpdateCat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user.id").(int)
	catRequest := models.CatInsertRequest{}

	id := r.PathValue("id")

	json.NewDecoder(r.Body).Decode(&catRequest)

	catRequest.UserId = userId

	err := validate.Struct(catRequest)
	helper.PanicIfError(err)

	catRequest.UserId = catRequest.UserId

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	cat, err := models.GetCatById(r.Context(), tx, id)
	if err != nil {
		panic(exception.NewNotFoundError("id is not found"))
	}

	if cat.UserId != userId {
		panic(exception.NewNotFoundError("id is not found"))
	}

	exist := models.CountCatInMatch(r.Context(), tx, id)

	if exist > 0 && catRequest.Sex != cat.Sex {
		panic(exception.NewBadRequestError("Cannot update sex when cat requested to match"))
	}

	if exist > 0 {
		_ = models.UpdateCatWithoutSex(r.Context(), tx, id, catRequest)
	} else {
		_ = models.UpdateCatWithSex(r.Context(), tx, id, catRequest)
	}

	wraper := helper.WebResponse{
		Message: "success",
		Data: map[string]interface{}{
			"id": id,
		},
	}

	helper.WriteToResponseBody(w, wraper, http.StatusOK)
}
