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
	var data interface{}
	catParam := models.CatParam{
		Id:            r.URL.Query().Get("id"),
		Owned:         r.URL.Query().Get("owned"),
		Email:         r.Header.Get("email"),
		AgeStr:        r.URL.Query().Get("ageInMonth"),
		HasMatchedStr: r.URL.Query().Get("hasMatched"),
		Race:          r.URL.Query().Get("race"),
		Search:        r.URL.Query().Get("search"),
		Limit:         r.URL.Query().Get("limit"),
		Offsset:       r.URL.Query().Get("offset"),
	}

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	if catParam.Id != "" {
		id, err := strconv.Atoi(catParam.Id)
		if err == nil {
			data, err = models.GetCatById(r.Context(), tx, id)
			if err != nil {
				data = nil
			}

		}
	} else {
		data = models.GetAllCat(r.Context(), tx, catParam)
	}

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

	if exist > 0 && catRequest.Sex != cat.Sex {
		panic(exception.NewBadRequestError("Cannot update sex when cat requested to match"))
	}

	if exist > 0 {
		err = models.UpdateCatWithSex(r.Context(), tx, id, catRequest)
	} else {
		err = models.UpdateCatWithoutSex(r.Context(), tx, id, catRequest)
	}

	wraper := helper.WebResponse{
		Message: "success",
		Data: map[string]interface{}{
			"id": fmt.Sprintf("%d", id),
		},
	}

	helper.WriteToResponseBody(w, wraper, http.StatusCreated)
}
