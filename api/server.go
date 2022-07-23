package api

import (
	db "blog/db/sqlc"
	"blog/token"
	"blog/util"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	store      db.Store
	config     util.Config
	tokenMaker token.Maker
}

func NewServer(store db.Store, config util.Config) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(AccessLog())
	router.POST("/auth/register", server.createAuth)
	router.POST("/auth/login", server.loginAuth)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/articles", server.CreateBlogArticle)
	authRoutes.POST("articles/update", server.UpdateBlogArticle)
	authRoutes.GET("/articles", server.listBlogAtricles)
	authRoutes.GET("/all_articles", server.listAllBlogAtricles)
	authRoutes.GET("/articles/:id", server.getBlogArticle)
	authRoutes.POST("/articles/delete", server.deleteArticle)

	authRoutes.POST("/tags", server.createBlogTag)
	authRoutes.POST("/tags/delete", server.deleteBlogTag)
	authRoutes.POST("tags/update", server.updateBlogTag)
	authRoutes.GET("/tags", server.listBlogTag)
	authRoutes.GET("/tags/:id", server.getBlogTag)

	server.router = router
}

//监听api请求
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
