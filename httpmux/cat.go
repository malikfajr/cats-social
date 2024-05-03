package httpmux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	catRequest := models.CatInsertRequest{}

	json.NewDecoder(r.Body).Decode(&catRequest)

	err := validate.Struct(catRequest)
	helper.PanicIfError(err)

	catRequest.UserEmail = r.Header.Get("email")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)
	id, date := models.SaveCat(r.Context(), tx, catRequest)

	wraper := helper.WebResponse{
		Message: "success",
		Data: map[string]interface{}{
			"id":        fmt.Sprintf("%d", id),
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
		Email:         r.Header.Get("email"),
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

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	data = models.GetAllCat(r.Context(), tx, catParam)

	wrapper := helper.WebResponse{
		Message: "success",
		Data:    data,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusOK)
}

func DestroyCat(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(exception.NewNotFoundError("id not found"))
	}

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	err = models.DestroyCat(r.Context(), tx, id, email)
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
	email := r.Header.Get("email")
	catRequest := models.CatInsertRequest{}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(exception.NewNotFoundError("id is not found"))
	}

	json.NewDecoder(r.Body).Decode(&catRequest)

	err = validate.Struct(catRequest)
	helper.PanicIfError(err)

	catRequest.UserEmail = r.Header.Get("email")

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	cat, err := models.GetCatById(r.Context(), tx, id)
	if err != nil {
		panic(exception.NewNotFoundError("id is not found"))
	}

	if cat.UserEmail != email {
		panic(exception.NewNotFoundError("id is not found"))
	}

	exist := models.CountCatInMatch(r.Context(), tx, idStr)
fmt.Println(catRequest.Sex, cat.Sex, exist)
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
			"id": fmt.Sprintf("%d", id),
		},
	}

	helper.WriteToResponseBody(w, wraper, http.StatusOK)
}
