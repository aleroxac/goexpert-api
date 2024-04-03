package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aleroxac/goexpert-api/internal/dto"
	"github.com/aleroxac/goexpert-api/internal/entity"
	"github.com/aleroxac/goexpert-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: userDB,
	}
}

// GetJWT godoc
//
//	@Summary		Get a user JWT
//	@Description	Get a user JWT
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request					body		dto.JWTInput	true	"user credentials"
//	@Success		200						{object}	dto.JWTOutput
//	@Failure		404						{object}	dto.Error
//	@Failure		401						{object}	dto.Error
//	@Failure		400						{object}	dto.Error
//	@Router			/users/generate_token 	[post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.JWTInput
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, tokenString, err := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	accessToken := dto.JWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
//
//	@Summary		Create 	user
//	@Description	Create 	user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.UserInput	true	"user request"
//	@Success		201
//	@Failure		500		{object}	dto.Error
//	@Failure		400		{object}	dto.Error
//	@Router			/users 	[post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.UserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	p, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.UserDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
