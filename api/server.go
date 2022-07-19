package api

import (
	db "blog/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  db.Store
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/articles", server.CreateBlogArticle)
	router.POST("articles/update", server.UpdateBlogArticle)
	router.GET("/articles", server.listBlogAtricles)
	router.GET("/articles/:id", server.getBlogArticle)
	router.POST("/articles/delete", server.deleteArticle)

	router.POST("/tags", server.createBlogTag)
	router.POST("/tags/delete", server.deleteBlogTag)
	router.POST("tags/update", server.updateBlogTag)
	router.GET("/tags", server.listBlogTag)
	router.GET("/tags/:id", server.getBlogTag)

	router.POST("/register", server.createAuth)
	router.POST("/login", server.loginAuth)

	server.router = router
}

//监听api请求
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
