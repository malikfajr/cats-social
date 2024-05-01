package httpmux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		AgeStr:        r.Header.Get("ageInMonth"),
		HasMatchedStr: r.Header.Get("hasMatched"),
		Race:          r.Header.Get("race"),
		Search:        r.Header.Get("search"),
		Limit:         r.Header.Get("limit"),
		Offsset:       r.Header.Get("offset"),
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
