package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/SraReaper/gofinance-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	CategoryID  int32     `json:"category_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Value       int32     `json:"value" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
}

// createAccount para criar uma conta
func (server *Server) createAccount(ctx *gin.Context) {
	var request createAccountRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	//Validação
	var categoryId = request.CategoryID
	var accountType = request.Type

	category, err := server.store.GetCategory(ctx, categoryId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	}
	var categoryTypeIsDifferentOfAccountType = category.Type != accountType
	if categoryTypeIsDifferentOfAccountType {
		ctx.JSON(http.StatusBadRequest, "Account type is different of Category type")
	} else {
		arg := db.CreateAccountParams{
			UserID:      request.UserID,
			CategoryID:  categoryId,
			Title:       request.Title,
			Type:        accountType,
			Description: request.Description,
			Value:       request.Value,
			Date:        request.Date,
		}

		account, err := server.store.CreateAccount(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		ctx.JSON(http.StatusOK, account)
	}
	//fim

}

type getAccountRequest struct {
	ID int32 `json:"id" binding:"required"`
}

// getAccount valida a URL e a conta
func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountReportsRequest struct {
	UserID int32  `uri:"user_id" binding:"required"`
	Type   string `uri:"type" binding:"required"`
}

func (server *Server) getAccountReports(ctx *gin.Context) {
	var request getAccountReportsRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountsReportsParams{
		UserID: request.UserID,
		Type:   request.Type,
	}

	sumReports, err := server.store.GetAccountsReports(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sumReports)
}

type getAccountGraphRequest struct {
	UserID int32  `uri:"user_id" binding:"required"`
	Type   string `uri:"type" binding:"required"`
}

func (server *Server) getAccountGraph(ctx *gin.Context) {
	var request getAccountGraphRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountsGraphParams{
		UserID: request.UserID,
		Type:   request.Type,
	}

	countGraph, err := server.store.GetAccountsGraph(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, countGraph)
}

type deleteAccountRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

// deleteAccount deleta a conta
func (server *Server) deleteAccount(ctx *gin.Context) {
	var request deleteAccountRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err = server.store.DeleteAccount(ctx, request.ID)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, true)
}

type updateAccountRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       int32  `json:"value"`
}

// updateAccount para atualizar uma conta
func (server *Server) updateAccount(ctx *gin.Context) {
	var request updateAccountRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateAccountParams{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
		Value:       request.Value,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	CategoryID  int32     `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// getAccount valida a URL e as contas
func (server *Server) getAccounts(ctx *gin.Context) {
	var request getAccountsRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountsParams{
		UserID:      request.UserID,
		Type:        request.Type,
		CategoryID:  request.CategoryID,
		Title:       request.Title,
		Description: request.Description,
		Date:        request.Date,
	}

	account, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
