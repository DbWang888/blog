package api

import (
	db "blog/db/sqlc"
	"blog/e"
	"blog/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedOn time.Time `json:"created_on"`
}

func getAuthResponse(auth db.BlogAuth) authResponse {
	return authResponse{
		UserID:    auth.ID,
		Username:  auth.Username.String,
		CreatedOn: auth.CreatedOn,
	}

}

type authResponseT struct {
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"created_on"`
}

func getAuthResponseT(auth db.BlogAuth) authResponseT {
	return authResponseT{
		UserID:    auth.ID,
		Username:  auth.Username.String,
		Password:  auth.Password.String,
		CreatedOn: auth.CreatedOn,
	}

}

type createAuthRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=6"`
}

type registerResponse struct {
	Auth authResponse `json:"auth"`
	Tag  tagResponse  `json:"tag"`
}

func (server *Server) createAuth(c *gin.Context) {
	var req createAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	arg := db.RegisterParams{
		Auth: db.BlogAuth{
			Username: util.NewSqlNullString(req.Username),
			Password: util.NewSqlNullString(hashPassword),
		},
	}

	registerResult, err := server.store.RegisterTX(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	var result registerResponse
	result.Auth = getAuthResponse(registerResult.Auth)
	result.Tag = getTagResponse(registerResult.Tag)

	c.JSON(http.StatusOK, e.GetSucessResult(result))

}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Auth        authResponse `json:"auth"`
	AccessToken string       `json:"access_token"`
}

func (server *Server) loginAuth(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	auth, err := server.store.GetAuthByUserName(c, util.NewSqlNullString(req.Username))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	fmt.Print(auth.Password.String)

	err = util.CheckHashPassword(req.Password, auth.Password.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR_AUTH, err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(auth.Username.String, int(auth.ID), server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	loginres := loginResponse{
		Auth:        getAuthResponse(auth),
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, e.GetSucessResult(loginres))

}
