package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
	"github.com/lib/pq"
)

type UserController struct {
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

// NewUserController creates new account controller
func NewUserController(store db.Store, tokenMaker token.Maker, config util.Config) UserController {
	return UserController{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

// alphanum: username should contian ASCII alphanumeric characters only
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// createUser godoc
// @Summary Create a User
// @Description Create User by json user params
// @Tags users
// @Accept  json
// @Produce  json
// @Param username body string true "user name"
// @Param password body string true "passward minLength(6)"
// @Param fullname body string true "full name"
// @Param email body string true "email"
// @Success 200 {object} userResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := c.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res := newUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

// loginUser godoc
// @Summary User Login
// @Description Login with username and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param username body string true "user name"
// @Param password body string true "passward minLength(6)"
// @Success 200 {object} loginUserResponse
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/login [post]
func (c *UserController) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	user, err := c.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	accessToken, err := c.tokenMaker.CreateToken(user.Username, c.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
