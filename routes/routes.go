package routes

import (
	"github.com/enehizy/go-cms/controllers"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(r *gin.RouterGroup) {
	createPageRoutes(r)
	createMediaRoutes(r)
	createPostRoutes(r)
}

func createPostRoutes(r *gin.RouterGroup) {
	page_route := r.Group("/posts")
	{
		page_route.GET("/", controllers.GetPosts)
		page_route.GET("/:id", controllers.GetPost)
		page_route.POST("/", controllers.CreatePost)
		page_route.PUT("/:id", controllers.UpdatePost)
		page_route.DELETE("/:id", controllers.DeletePost)
	}
}

func createPageRoutes(r *gin.RouterGroup) {
	page_route := r.Group("/pages")
	{
		page_route.GET("/", controllers.GetPages)
		page_route.GET("/:id", controllers.GetPage)
		page_route.POST("/", controllers.CreatePage)
		page_route.PUT("/:id", controllers.UpdatePage)
		page_route.DELETE("/:id", controllers.DeletePage)
	}
}

func createMediaRoutes(r *gin.RouterGroup) {
	page_route := r.Group("/media")
	{
		page_route.GET("/", controllers.GetMedias)
		page_route.GET("/:id", controllers.GetMedia)
		page_route.POST("/", controllers.CreateMedia)
		page_route.PUT("/:id", controllers.UpdateMedia)
		page_route.DELETE("/:id", controllers.DeleteMedia)
	}
}
