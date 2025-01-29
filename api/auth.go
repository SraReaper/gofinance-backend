package api

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (server *Server) login(ctx *gin.Context) {
	var request loginRequest
	err := ctx.ShouldBindJSON(&request)
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

	hashedInput := sha512.Sum512([]byte(request.Password))
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	preparedPassword := string(trimmedHash)
	plainTextInBytes := []byte(preparedPassword)
	hashTextInBytes := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashTextInBytes, plainTextInBytes)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: request.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtSignedKey = []byte("secret_key")
	generatedTokenToString, err := generatedToken.SignedString(jwtSignedKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, generatedTokenToString)
}
