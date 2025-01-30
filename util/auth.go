package util

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ValidateToken(ctx *gin.Context, token string) error {
	claims := &Claims{}
	var jwtSignedKey = []byte("secret_key")
	tokenParse, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtSignedKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, err)
			return err
		}
		ctx.JSON(http.StatusBadRequest, err)
		return err
	}

	if !tokenParse.Valid {
		ctx.JSON(http.StatusUnauthorized, "Token is invalid")
		return nil
	}

	ctx.Next()
	return nil
}

func GetTokenInHeaderAndVerify(ctx *gin.Context) error {
	authorizationHeaderKey := ctx.GetHeader("authorization")
	fields := strings.Fields(authorizationHeaderKey)
	if len(fields) < 2 || fields[0] != "Bearer" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
		return errors.New("invalid authorization format")
	}
	tokenToValidate := fields[1]
	errOnValiteToken := ValidateToken(ctx, tokenToValidate)
	if errOnValiteToken != nil {
		return errOnValiteToken
	}
	return nil
}
