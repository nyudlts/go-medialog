package medialog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
)

func LoadRoutes(router *gin.Engine) {

	//Main Index
	router.GET("", func(c *gin.Context) { controllers.GetIndex(c) })
	router.GET("/test", func(c *gin.Context) { Test(c) })
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
	entryRoutes.GET("/:id/previous", func(c *gin.Context) { controllers.GetPreviousEntry(c) })
	entryRoutes.GET("/:id/next", func(c *gin.Context) { controllers.GetNextEntry(c) })

	//Users Group
	userRoutes := router.Group("/users")
	userRoutes.GET("", func(c *gin.Context) { controllers.GetUsers(c) })
	userRoutes.GET("/new", func(c *gin.Context) { controllers.NewUser(c) })
	userRoutes.POST("/create", func(c *gin.Context) { controllers.CreateUser(c) })
	userRoutes.GET("/login", func(c *gin.Context) { controllers.LoginUser(c) })
	userRoutes.POST("/authenticate", func(c *gin.Context) { controllers.AuthenticateUser(c) })
	userRoutes.GET("/:id/reset_password", func(c *gin.Context) { controllers.ResetUserPassword(c) })
	userRoutes.POST("/:id/reset_password", func(c *gin.Context) { controllers.ResetPassword(c) })
	userRoutes.GET("/:id/deactivate", func(c *gin.Context) { controllers.DeactivateUser(c) })
	userRoutes.GET("/:id/reactivate", func(c *gin.Context) { controllers.ReactivateUser(c) })
	userRoutes.GET("/:id/make_admin", func(c *gin.Context) { controllers.MakeUserAdmin(c) })
	userRoutes.GET("/:id/remove_admin", func(c *gin.Context) { controllers.RemoveUserAdmin(c) })

	//Report Group
	reportRoutes := router.Group("/reports")
	reportRoutes.GET("", func(c *gin.Context) { controllers.ReportsIndex(c) })
	reportRoutes.GET("/year", func(c *gin.Context) { controllers.ReportYear(c) })
	reportRoutes.POST("/range", func(c *gin.Context) { controllers.ReportRange(c) })

	//Search Group
	searchRoutes := router.Group("/search")
	searchRoutes.POST("", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, "Not Implemented") })

	//Session Group
	sessionRoutes := router.Group("/sessions")
	sessionRoutes.GET("/logout", func(c *gin.Context) { controllers.LogoutUser(c) })
	sessionRoutes.GET("/dump", func(c *gin.Context) { controllers.DumpSession(c) })
}

func Test(c *gin.Context) {
	c.JSON(200, "This is a Test")
}
