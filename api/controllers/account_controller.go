package controllers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hhow09/simple_bank/constants"
	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
	"github.com/lib/pq"
)

type AccountController struct {
	store db.Store
}

// AccountController creates new account controller
func NewAccountController(store db.Store) AccountController {
	return AccountController{
		store: store,
	}
}

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// CreateAccount godoc
// @Summary Create Account
// @Description create account by a already-login user
// @Tags accounts
// @Accept  json
// @Produce  json
// @Security authorization
// @Param currency body string true "currency"
// @Success 200 {object} CreateAccountRequest
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Router /accounts [post]
func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet(constants.AuthPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := c.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount godoc
// @Summary get Account
// @Description get account by account id
// @Tags accounts
// @Accept  json
// @Produce  json
// @Security authorization
// @Param id path integer true "Account ID"
// @Success 200 {object} db.Account
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Router /accounts/:id [get]
func (c *AccountController) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := c.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(constants.AuthPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belongs to the authenticated user.")
		ctx.JSON(http.StatusUnauthorized, err)
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts godoc
// @Summary list Account
// @Description list account under current user
// @Tags accounts
// @Accept  json
// @Produce  json
// @Security authorization
// @Param page_id query int true "page id minimum(1)"
// @Param page_size query int true "page minimum(5) maximum(10)"
// @Success 200 {object} []db.Account
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Router /accounts [get]
func (c *AccountController) ListAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet(constants.AuthPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := c.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
