package api

import (
	db "blog/db/sqlc"
	"blog/e"
	"blog/util"
	"errors"
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

	arg := db.RegisterParams{
		Auth: db.BlogAuth{
			Username: util.NewSqlNullString(req.Username),
			Password: util.NewSqlNullString(req.Password),
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

	if req.Password != auth.Password.String {
		err = errors.New("密码错误")
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR_AUTH_CHECK_TOKEN_FAIL, err))
		return
	}

	c.JSON(http.StatusOK, e.GetSucessResult(getAuthResponse(auth)))

}
