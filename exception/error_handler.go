package exception

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/cats-social/helper"
)

func RecoverWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if validationError(w, r, err) {
					return
				}

				if badRequestError(w, r, err) {
					return
				}

				if notFoundError(w, r, err) {
					return
				}

				if confilctError(w, r, err) {
					return
				}

				internalServerError(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func confilctError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(ConflictError)
	if ok {
		wrapper := helper.WebResponse{
			Message: exception.Error,
			Data:    nil,
		}
		helper.WriteToResponseBody(writer, wrapper, http.StatusConflict)
		return true
	} else {
		return false
	}
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		wrapper := helper.WebResponse{
			Message: exception.Error,
			Data:    nil,
		}
		helper.WriteToResponseBody(writer, wrapper, http.StatusNotFound)
		return true
	} else {
		return false
	}
}

func badRequestError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		wrapper := helper.WebResponse{
			Message: exception.Error,
			Data:    nil,
		}
		helper.WriteToResponseBody(writer, wrapper, http.StatusBadRequest)
		return true
	} else {
		return false
	}
}
func validationError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		webResponse := helper.WebResponse{
			Message: exception.Error(),
			Data:    nil,
		}

		helper.WriteToResponseBody(writer, webResponse, http.StatusBadRequest)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	fmt.Println(err)

	webResponse := helper.WebResponse{
		Message: "Internal server error",
		Data:    nil,
	}

	helper.WriteToResponseBody(writer, webResponse, http.StatusInternalServerError)
}
