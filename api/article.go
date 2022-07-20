package api

import (
	db "blog/db/sqlc"
	"blog/e"
	"blog/token"
	"blog/util"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type articleResponse struct {
	ID         int32     `json:"id"`
	TagID      int32     `json:"tag_id"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	Content    string    `json:"content"`
	CreatedOn  time.Time `json:"created_on"`
	CreatedBy  string    `json:"created_by"`
	ModifiedOn time.Time `json:"modified_on"`
	ModifiedBy string    `json:"modified_by"`
	DeletedOn  time.Time `json:"deleted_on"`
	State      int32     `json:"state"`
}

func getArticleRes(article db.BlogArticle) articleResponse {
	return articleResponse{
		ID:         article.ID,
		TagID:      article.TagID.Int32,
		Title:      article.Title.String,
		Desc:       article.Desc.String,
		Content:    article.Content.String,
		CreatedOn:  article.CreatedOn,
		CreatedBy:  article.CreatedBy.String,
		ModifiedOn: article.ModifiedOn,
		ModifiedBy: article.ModifiedBy.String,
		DeletedOn:  article.DeletedOn,
		State:      article.State.Int32,
	}
}

//接收创建文章接口 前端post参数
type createBlogArticleRequest struct {
	TagID   int32  `json:"tag_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Desc    string `json:"desc" binding:"required"`
	Content string `json:"content" binding:"required"`
}

//创建文章接口
func (server *Server) CreateBlogArticle(c *gin.Context) {
	var req createBlogArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.INVALID_PARAMS, err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateBlogArticleParams{
		TagID:     util.NewSqlNullInt32(req.TagID),
		Title:     util.NewSqlNullString(req.Title),
		CreatedBy: util.NewSqlNullString(authPayload.UserName),
		Desc:      util.NewSqlNullString(req.Desc),
		Content:   util.NewSqlNullString(req.Content),
	}

	createdResult, err := server.store.CreateBlogArticle(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR_NOT_EXIST_ARTICLE, err))
		return
	}

	createdArticleID, err := createdResult.LastInsertId()
	if err != nil {
		c.JSON(http.StatusBadGateway, e.GetErrResult(e.ERROR, err))
		return
	}

	article, err := server.store.GetBlogArticles(c, int32(createdArticleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}
	res := getArticleRes(article)

	result := e.GetSucessResult(res)

	c.JSON(http.StatusOK, result)
}

//修改文章前端参数
type updateArticleRequest struct {
	TagID   int32  `json:"tag_id" binding:"required"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

//修改文章实现
//通过查询需要修改的内容对比，是否为修改内容，再将对应的内容进行修改
func (server *Server) UpdateBlogArticle(c *gin.Context) {
	var req updateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.ERROR, err))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	sqlarticle, err := server.store.GetBlogArticles(c, int32(authPayload.UserID))
	if err != nil {
		log.Fatal("update article -- getblogartice failed", err)
		c.JSON(http.StatusBadGateway, e.GetErrResult(e.ERROR, err))
		return
	}

	article := getArticleRes(sqlarticle)

	if req.Title == "" {
		req.Title = article.Title
	}

	if req.Content == "" {
		req.Content = article.Content
	}

	if req.Desc == "" {
		req.Desc = article.Desc
	}

	arg := db.UpdateBlogArticleParams{
		TagID:      util.NewSqlNullInt32(req.TagID),
		Title:      util.NewSqlNullString(req.Title),
		Desc:       util.NewSqlNullString(req.Desc),
		Content:    util.NewSqlNullString(req.Content),
		ModifiedBy: util.NewSqlNullString(authPayload.UserName),
		ModifiedOn: time.Now(),
		ID:         int32(authPayload.UserID),
	}

	_, err = server.store.UpdateBlogArticle(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	sqlarticle, err = server.store.GetBlogArticles(c, int32(authPayload.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.GetErrResult(e.ERROR, err))
		return
	}

	article = getArticleRes(sqlarticle)
	result := e.GetSucessResult(article)
	c.JSON(http.StatusOK, result)

}

//展示用户本人的文章列表
type listArticleRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listBlogAtricles(c *gin.Context) {
	var req listArticleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.ERROR, err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListBlogAtriclesParams{
		CreatedBy: util.NewSqlNullString(authPayload.UserName),
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	Sqlarticles, err := server.store.ListBlogAtricles(c, arg)
	log.Printf("%v", err)
	if err != nil {
		c.JSON(http.StatusNotFound, e.GetErrResult(e.ERROR, err))
		return
	}

	len_Sqlarticles := len(Sqlarticles)

	articles := make([]articleResponse, len_Sqlarticles)
	for i := 0; i < len_Sqlarticles; i++ {
		articles[i] = getArticleRes(Sqlarticles[i])

	}

	result := e.GetSucessResult(articles)
	c.JSON(http.StatusOK, result)

}

//展示所有人的文章，作为首页
func (server *Server) listAllBlogAtricles(c *gin.Context) {
	var req listArticleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.ERROR, err))
		return
	}

	arg := db.ListAllArticlesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	Sqlarticles, err := server.store.ListAllArticles(c, arg)
	log.Printf("%v", err)
	if err != nil {
		c.JSON(http.StatusNotFound, e.GetErrResult(e.ERROR, err))
		return
	}

	len_Sqlarticles := len(Sqlarticles)

	articles := make([]articleResponse, len_Sqlarticles)
	for i := 0; i < len_Sqlarticles; i++ {
		articles[i] = getArticleRes(Sqlarticles[i])

	}

	result := e.GetSucessResult(articles)
	c.JSON(http.StatusOK, result)
}

//返回指定文章内容
type getArticleRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBlogArticle(c *gin.Context) {

	var req getArticleRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusForbidden, e.GetErrResult(e.ERROR, err))
		return
	}

	article, err := server.store.GetBlogArticles(c, req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, e.GetErrResult(e.ERROR, err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	res := getArticleRes(article)

	if res.CreatedBy != authPayload.UserName {
		err := errors.New("不是本人创建，无权查看")
		c.JSON(http.StatusUnauthorized, e.GetErrResult(e.ERROR_AUTH_TOKEN, err))
	}

	result := e.GetSucessResult(res)

	c.JSON(http.StatusOK, result)

}

//删除指定文章内容
type deleteArticleRequest struct {
	ID int32 `form:"id" binding:"required,min=1"`
}

func (server *Server) deleteArticle(c *gin.Context) {
	var req deleteArticleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.GetErrResult(e.ERROR, err))
		return
	}

	arg := db.DeleteArticleParams{
		DeletedOn: time.Now(),
		ID:        req.ID,
	}

	_, err := server.store.DeleteArticle(c, arg)

	if err != nil {
		c.JSON(http.StatusForbidden, e.GetErrResult(e.ERROR, err))
		return
	}

	result := e.GetSucessResult(nil)
	c.JSON(http.StatusOK, result)

}
