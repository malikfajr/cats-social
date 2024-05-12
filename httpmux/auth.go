package httpmux

import (
	"encoding/json"
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

	json.NewDecoder(r.Body).Decode(&credential)

	err := validate.Struct(credential)
	helper.PanicIfError(err)

	db := models.GetDb()

	user, err := models.GetUserByEmail(r.Context(), db, credential.Email)
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
		"accessToken": generateToken(user.Id, user.Email, user.Name),
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
	json.NewDecoder(r.Body).Decode(&user)

	// validation json
	err := validate.Struct(user)
	helper.PanicIfError(err)

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), config.Env.BCRYPT_SALT)
	helper.PanicIfError(err)

	user.Password = string(hashPassword)

	tx := models.StartTx()
	defer helper.CommitOrRollback(tx)

	newUser := models.SaveUser(r.Context(), tx, user)

	data := map[string]string{
		"email":       newUser.Email,
		"name":        newUser.Name,
		"accessToken": generateToken(newUser.Id, newUser.Email, newUser.Name),
	}

	wrapper := helper.WebResponse{
		Message: "User registered successfully",
		Data:    data,
	}

	helper.WriteToResponseBody(w, wrapper, http.StatusCreated)
}

func generateToken(userId int, email string, name string) string {
	myClaims := config.CustomJWTClaim{
		ID:    userId,
		Email: email,
		Name:  name,
		Exp:   jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)

	ss, err := token.SignedString([]byte(config.Env.JWT_SECRET))
	helper.PanicIfError(err)

	return ss
}
