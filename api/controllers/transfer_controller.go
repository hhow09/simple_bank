package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hhow09/simple_bank/constants"
	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
)

type TransferController struct {
	store  db.Store
	config util.Config
}

func NewTransferController(store db.Store, tokenMaker token.Maker, config util.Config) TransferController {
	return TransferController{
		store:  store,
		config: config,
	}
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// CreateTransfer godoc
// @Summary Create Transfer
// @Description Create transfer from from_account_id to to_account_id which has same currency
// @Tags transfers
// @Accept  json
// @Produce  json
// @Security authorization
// @Param from_account_id body integer true "from_account_id"
// @Param to_account_id body integer true "to_account_id"
// @Param amount body integer true "amount"
// @Param currency body string true "currency"
// @Success 200 {object} db.TransferTxResult
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Router /transfers [post]
func (c *TransferController) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	fromAccount, valid := c.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(constants.AuthPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belongs to the current user")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	_, valid = c.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := c.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *TransferController) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := c.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return account, false
	}

	return account, true
}
