package medialog

import (
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
)

func LoadRoutes(router *gin.Engine) {

	//Main Index
	router.GET("", func(c *gin.Context) { controllers.GetIndex(c) })

	//Accessions Group
	accessionRoutes := router.Group("/accessions")
	accessionRoutes.GET("", func(c *gin.Context) { controllers.GetAccessions(c) })
	accessionRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetAccession(c) })

	//Repository Group
	repositoryRoutes := router.Group("/repositories")
	repositoryRoutes.GET("", func(c *gin.Context) { controllers.GetRepositories(c) })
	repositoryRoutes.GET("/new", func(c *gin.Context) { controllers.NewRepository(c) })
	repositoryRoutes.POST("", func(c *gin.Context) { controllers.CreateRepository(c) })
	repositoryRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetRepository(c) })

	//Resources Group
	resourceRoutes := router.Group("/resources")
	resourceRoutes.GET("", func(c *gin.Context) { controllers.GetResources(c) })
	resourceRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetResource(c) })

	//Entries Group
	entryRoutes := router.Group("/entries")
	entryRoutes.GET("", func(c *gin.Context) { controllers.GetEntries(c) })
	entryRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetEntry(c) })

}
