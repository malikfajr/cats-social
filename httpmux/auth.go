package httpmux

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/malikfajr/cats-social/config"
	"github.com/malikfajr/cats-social/exception"
	"github.com/malikfajr/cats-social/helper"
	"github.com/malikfajr/cats-social/models"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=5,max=15"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	credential := Credential{}

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err == io.EOF {
		panic(exception.NewBadRequestError("missing data"))
	}
	helper.PanicIfError(err)

	err = validate.Struct(credential)
	helper.PanicIfError(err)

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	user, err := models.GetUserByEmail(r.Context(), tx, credential.Email)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
	if err != nil {
		panic(exception.NewBadRequestError("Password wrong"))
	}

	data := map[string]string{
		"email":       user.Email,
		"name":        user.Name,
		"accessToken": generateToken(user.Email),
	}

	wrapper := helper.WebResponse{
		Message: "User logged successfully",
		Data:    data,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusOK)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	// parsing body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err == io.EOF {
		panic(exception.NewBadRequestError("missing data"))
	}
	helper.PanicIfError(err)

	// validation json
	err = validate.Struct(user)
	helper.PanicIfError(err)

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), config.Env.BCRYPT_SALT)
	helper.PanicIfError(err)

	user.Password = string(hashPassword)

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	_, err = models.GetUserByEmail(r.Context(), tx, user.Email)
	if err == nil {
		panic(exception.NewConflictError("Email has taken"))
	}

	newUser := models.SaveUser(r.Context(), tx, user)

	data := map[string]string{
		"email":       newUser.Email,
		"name":        newUser.Name,
		"accessToken": generateToken(newUser.Name),
	}

	wrapper := helper.WebResponse{
		Message: "User registered successfully",
		Data:    data,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusCreated)
}

func generateToken(email string) string {
	myClaims := config.CustomJWTClaim{
		Email: email,
		Exp:   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)

	ss, err := token.SignedString([]byte(config.Env.JWT_SECRET))
	helper.PanicIfError(err)

	return ss
}

func Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Protected" + r.Header.Get("email")))
}
