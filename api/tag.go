package api

import (
	db "blog/db/sqlc"
	"blog/e"
	"blog/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type tagResponse struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	CreatedOn  time.Time `json:"created_on"`
	CreatedBy  string    `json:"created_by"`
	ModifiedOn time.Time `json:"modified_on"`
	ModifiedBy string    `json:"modified_by"`
	DeletedOn  time.Time `json:"deleted_on"`
	State      int32     `json:"state"`
}

//格式化返回数据
func getTagResponse(blogtag db.BlogTag) tagResponse {
	return tagResponse{
		ID:         blogtag.ID,
		Name:       blogtag.Name.String,
		CreatedOn:  blogtag.CreatedOn,
		CreatedBy:  blogtag.CreatedBy.String,
		ModifiedOn: blogtag.ModifiedOn,
		ModifiedBy: blogtag.ModifiedBy.String,
		DeletedOn:  blogtag.DeletedOn,
		State:      blogtag.State.Int32,
	}
}

//接收创建tag参数
type createTagRequest struct {
	Name      string `json:"name" binding:"required"`
	CreatedBy string `json:"created_by" binding:"required"`
}

//实现创建tag接口
func (server *Server) createBlogTag(c *gin.Context) {
	var req createTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	arg := db.CreateBlogTagParams{
		Name:      util.NewSqlNullString(req.Name),
		CreatedBy: util.NewSqlNullString(req.CreatedBy),
	}

	createdResult, err := server.store.CreateBlogTag(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	createdID, err := createdResult.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	createdTag, err := server.store.GetBlogTag(c, int32(createdID))
	if err != nil {
		c.JSON(http.StatusNotFound, e.GetErrResult(e.ERROR_NOT_EXIST_TAG, err))
		return
	}

	result := getTagResponse(createdTag)

	c.JSON(http.StatusOK, e.GetSucessResult(result))

}

type deleteTagRequest struct {
	ID int32 `form:"id" binding:"required,min=1"`
}

func (server *Server) deleteBlogTag(c *gin.Context) {
	var req deleteTagRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	arg := db.DeleteBlogTagParams{
		ID:        req.ID,
		DeletedOn: time.Now(),
	}

	_, err := server.store.DeleteBlogTag(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	c.JSON(http.StatusOK, e.GetSucessResult(nil))
}

type updateTagRequest struct {
	Name       string `json:"name" binding:"required"`
	ModifiedBy string `json:"modified_by" binding:"required"`
	ID         int32  `json:"id" binding:"required,min=1"`
}

func (server *Server) updateBlogTag(c *gin.Context) {
	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	arg := db.UpdateBlogTagParams{
		Name:       util.NewSqlNullString(req.Name),
		ModifiedBy: util.NewSqlNullString(req.ModifiedBy),
		ModifiedOn: time.Now(),
		ID:         req.ID,
	}

	_, err := server.store.UpdateBlogTag(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	c.JSON(http.StatusOK, e.GetSucessResult(nil))
}

type listTagRequest struct {
	CreatedBy string `form:"created_by" binding:"required"`
	PageID    int32  `form:"page_id" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listBlogTag(c *gin.Context) {
	var req listTagRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	arg := db.ListBlogTagParams{
		CreatedBy: util.NewSqlNullString(req.CreatedBy),
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	SqlTags, err := server.store.ListBlogTag(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	sqlTagsLen := len(SqlTags)

	tags := make([]tagResponse, sqlTagsLen)

	for i, v := range SqlTags {
		tags[i] = getTagResponse(v)
	}

	c.JSON(http.StatusOK, e.GetSucessResult(tags))
}

type getTagRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBlogTag(c *gin.Context) {
	var req getTagRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusForbidden, e.GetErrResult(e.ERROR, err))
		return
	}

	tag, err := server.store.GetBlogTag(c, req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, e.GetErrResult(e.ERROR, err))
		return
	}

	res := getTagResponse(tag)
	result := e.GetSucessResult(res)

	c.JSON(http.StatusOK, result)
}
