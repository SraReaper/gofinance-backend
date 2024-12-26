package api

import (
	"database/sql"
	"net/http"

	db "github.com/SraReaper/gofinance-backend/db/sqlc"
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
	var request deleteCategoryRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err = server.store.DeleteCategory(ctx, request.ID)
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
	var request updateCategoryRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateCategoryParams{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
	}

	category, err := server.store.UpdateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, category)
}

type getCategoriesRequest struct {
	UserID      int32  `json:"user_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// getCategories valida a URL e a categorias
func (server *Server) getCategories(ctx *gin.Context) {
	var request getCategoriesRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var categories []db.Category
	var parametersHasUserIdAndType = request.UserID > 0 && len(request.Type) > 0

	filterAsByUserIdAndType := len(request.Description) == 0 && len(request.Title) == 0 && parametersHasUserIdAndType
	if filterAsByUserIdAndType {
		arg := db.GetCategoriesByUserIdAndTypeParams{
			UserID: request.UserID,
			Type:   request.Type,
		}

		categoriesByUserIdAndType, err := server.store.GetCategoriesByUserIdAndType(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByUserIdAndType

	}

	filterAsByUserIdAndTypeAndDescription := len(request.Title) == 0 && len(request.Description) > 0 && parametersHasUserIdAndType
	if filterAsByUserIdAndTypeAndDescription {
		arg := db.GetCategoriesByUserIdAndTypeAndDescriptionParams{
			UserID:      request.UserID,
			Type:        request.Type,
			Description: request.Description,
		}

		categoriesByUserIdAndTypeAndDescription, err := server.store.GetCategoriesByUserIdAndTypeAndDescription(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByUserIdAndTypeAndDescription
	}

	filterAsByUserIdAndTypeAndTitle := len(request.Title) > 0 && len(request.Description) == 0 && parametersHasUserIdAndType
	if filterAsByUserIdAndTypeAndTitle {
		arg := db.GetCategoriesByUserIdAndTypeAndTitleParams{
			UserID: request.UserID,
			Type:   request.Type,
			Title:  request.Title,
		}

		categoriesByUserIdAndTypeAndTitle, err := server.store.GetCategoriesByUserIdAndTypeAndTitle(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesByUserIdAndTypeAndTitle
	}

	filterAsAllParameters := len(request.Title) > 0 && len(request.Description) > 0 && parametersHasUserIdAndType
	if filterAsAllParameters {
		arg := db.GetCategoriesParams{
			UserID:      request.UserID,
			Type:        request.Type,
			Title:       request.Title,
			Description: request.Description,
		}

		categoriesAllFilters, err := server.store.GetCategories(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		categories = categoriesAllFilters
	}

	ctx.JSON(http.StatusOK, categories)
}
