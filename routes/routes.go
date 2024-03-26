package medialog

import (
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
)

func LoadRoutes(router *gin.Engine) {

	//Main Index
	router.GET("", func(c *gin.Context) { c.JSON(200, "Hello") })

	//Accessions Group
	accessionsRoutes := router.Group("/accessions")
	accessionsRoutes.GET("", func(c *gin.Context) { controllers.GetAccessions(c) })
	accessionsRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetAccession(c) })

	//Repository Group
	repositoryRoutes := router.Group("/repositories")
	repositoryRoutes.GET("", func(c *gin.Context) { controllers.GetRepositories(c) })
	repositoryRoutes.GET("/new", func(c *gin.Context) { controllers.NewRepository(c) })
	repositoryRoutes.POST("", func(c *gin.Context) { controllers.CreateRepository(c) })
	repositoryRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetRepository(c) })

}
