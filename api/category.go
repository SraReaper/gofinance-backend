package api

import (
	"database/sql"
	"net/http"

	db "github.com/SraReaper/gofinance-backend/db/sqlc"
	"github.com/SraReaper/gofinance-backend/util"
	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	UserID      int32  `json:"user_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// createCategory para criar um usuário
func (server *Server) createCategory(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}

	var request createCategoryRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreateCategoryParams{
		UserID:      request.UserID,
		Title:       request.Title,
		Type:        request.Type,
		Description: request.Description,
	}

	user, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, user)
}

type getCategoryRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

// getCategory valida a URL e a categoria
func (server *Server) getCategory(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}

	var request getCategoryRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetCategory(ctx, request.ID)
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

type deleteCategoryRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

// deleteCategory deleta a categoria
func (server *Server) deleteCategory(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}

	var request deleteCategoryRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err = server.store.DeleteCategories(ctx, request.ID)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, true)
}

type updateCategoryRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// updateCategory para atualizar um usuário
func (server *Server) updateCategory(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var request updateCategoryRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateCategoriesParams{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
	}

	category, err := server.store.UpdateCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, category)
}

type getCategoriesRequest struct {
	UserID      int32  `form:"user_id" json:"user_id" binding:"required"`
	Type        string `form:"type" json:"type" binding:"required"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

// getCategories valida a URL e a categorias
func (server *Server) getCategories(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}

	var request getCategoriesRequest
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetCategoriesParams{
		UserID:      request.UserID,
		Type:        request.Type,
		Title:       request.Title,
		Description: request.Description,
	}

	categories, err := server.store.GetCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, categories)
}
