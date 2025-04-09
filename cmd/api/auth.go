package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gad/social/internal/mailer"
	"github.com/gad/social/internal/store"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,max=50"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	payload := RegisterUserPayload{}
	if err := readJson(app, w, r, &payload); err != nil {

		app.badRequestResponse(w, r, err, false)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Username); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	plainToken, err := uuid.NewV4()
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
	

	plainTokenS := plainToken.String()
	app.logger.Debugf("plain token of newly created %s", plainTokenS)
	hash:= sha256.Sum256([]byte(plainTokenS))
	hashToken := hex.EncodeToString(hash[:])

	err = app.store.Users.RegisterNew(r.Context(), user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrorDuplicateEmail:
			app.badRequestResponse(w, r, err, true)
		case store.ErrorDuplicateµUsername:
			app.badRequestResponse(w, r, err, true)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	vars := struct{
		Username string
		ActivationLink string
	} {
		Username: user.Username,
		ActivationLink: fmt.Sprintf("%s/activate/%s",app.config.frontendURL,plainTokenS),
	}

	isProdEnv := app.config.env == "PRODUCTION"
	err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv )
	if err != nil{
		app.logger.Errorw("Error sending confirmation email : ", err)

		// rollback user creation

		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("Error sending confirmation email : ", err)
			
		}
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
type SetUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
// setTokenHandler godoc
//
//	@Summary		Create a token
//	@Description	Create a token for a given user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		SetUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string				"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) setTokenHandler (w http.ResponseWriter, r *http.Request) {

		// parse r payload  -> credentials + process errors
		payload := SetUserTokenPayload{}
		if err := readJson(app, w, r, &payload); err != nil {
	
			app.badRequestResponse(w, r, err, false)
			return
		}
	
		if err := validate.Struct(payload); err != nil {
			app.badRequestResponse(w, r, err, true)
			return
		}
		// fetch user by email
		user, err := app.store.Users.GetUserByEmail(r.Context(),payload.Email)

		if err != nil {
			switch err {
				case store.ErrorNotFound:
					app.AuthError(w, r, err)
				default:
					app.internalServerErrorResponse(w, r, err)
			}
			return
		}

		if err = user.Password.Compare(payload.Password); err!= nil {
			app.AuthError(w, r, err)
			return
		}

		claims := jwt.MapClaims{
			"sub" : user.ID,
			"exp" : time.Now().Add(app.config.auth.token.exp).Unix(),
			"iat" : time.Now().Unix(),
			"nbf" : time.Now().Unix(),
			"iss" : app.config.auth.token.issuer,
			"aud" : app.config.auth.token.issuer,


		}
		token, err := app.authenticator.GenerateToken(claims)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
		}

		// add claim

		// generate token 
		// sign token
		// send token
		if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
			app.internalServerErrorResponse(w, r, err)
		}
}
