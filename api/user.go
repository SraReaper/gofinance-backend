package api

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"net/http"

	db "github.com/SraReaper/gofinance-backend/db/sqlc"
	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// createUser para criar um usuário
func (server *Server) createUser(ctx *gin.Context) {
	var request createUserRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	hashedInput := sha512.Sum512([]byte(request.Password))
	trimmedHash := bytes.Trim(hashedInput[:], "/x00")
	preparedPassword := string(trimmedHash)
	passwordHashInBytes, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	var passwordHashed = string(passwordHashInBytes)

	arg := db.CreateUserParams{
		Username: request.Username,
		Password: passwordHashed,
		Email:    request.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

// getUser valida a URL e o usuario
func (server *Server) getUser(ctx *gin.Context) {
	var request getUserRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUser(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserByIdRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

// getUserById valida o id do usuario
func (server *Server) getUserById(ctx *gin.Context) {
	var request getUserByIdRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUserById(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
