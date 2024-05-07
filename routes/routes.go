package medialog

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
)

func LoadRoutes(router *gin.Engine) {
	// defaults
	router.NoRoute(func(c *gin.Context) { c.JSON(404, "NOT FOUND") })
	router.NoMethod(func(c *gin.Context) { c.JSON(405, "NO METHOD") })
	router.GET("/401", func(c *gin.Context) {
		session := sessions.Default(c)
		session.AddFlash("Please authenticate to access this service", "WARNING")
		c.HTML(401, "error.html", gin.H{"flash": session.Flashes("WARNING")})
		session.Save()
	})

	//Main Index
	router.GET("", func(c *gin.Context) { controllers.GetIndex(c) })
	router.GET("/test", func(c *gin.Context) { Test(c) })

	//Accessions Group
	accessionsRoutes := router.Group("/accessions")
	accessionsRoutes.GET("", func(c *gin.Context) { controllers.GetAccessions(c) })
	accessionsRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetAccession(c) })
	accessionsRoutes.GET("/:id/slew", func(c *gin.Context) { controllers.SlewAccession(c) })
	accessionsRoutes.POST("slew", func(c *gin.Context) { controllers.CreateAccessionSlew(c) })

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
	entryRoutes.GET("/new", func(c *gin.Context) { controllers.NewEntry(c) })
	entryRoutes.POST("", func(c *gin.Context) { controllers.CreateEntry(c) })
	entryRoutes.GET("/:id/edit", func(c *gin.Context) { controllers.EditEntry(c) })
	entryRoutes.POST("/:id/update", func(c *gin.Context) { controllers.UpdateEntry(c) })
	entryRoutes.GET("/:id/delete", func(c *gin.Context) { controllers.DeleteEntry(c) })
	entryRoutes.GET("/:id/show", func(c *gin.Context) { controllers.GetEntry(c) })
	entryRoutes.GET("/:id/previous", func(c *gin.Context) { controllers.GetPreviousEntry(c) })
	entryRoutes.GET("/:id/next", func(c *gin.Context) { controllers.GetNextEntry(c) })
	entryRoutes.GET("/:id/clone", func(c *gin.Context) { controllers.CloneEntry(c) })

	//Users Group
	userRoutes := router.Group("/users")
	userRoutes.GET("", func(c *gin.Context) { controllers.GetUsers(c) })
	userRoutes.GET("new", func(c *gin.Context) { controllers.NewUser(c) })
	userRoutes.POST("create", func(c *gin.Context) { controllers.CreateUser(c) })
	userRoutes.GET("login", func(c *gin.Context) { controllers.LoginUser(c) })
	userRoutes.GET("logout", func(c *gin.Context) { controllers.LogoutUser(c) })
	userRoutes.POST("authenticate", func(c *gin.Context) { controllers.AuthenticateUser(c) })
	userRoutes.GET(":id/reset_password", func(c *gin.Context) { controllers.ResetUserPassword(c) })
	userRoutes.POST(":id/reset_password", func(c *gin.Context) { controllers.ResetPassword(c) })
	userRoutes.GET(":id/deactivate", func(c *gin.Context) { controllers.DeactivateUser(c) })
	userRoutes.GET(":id/reactivate", func(c *gin.Context) { controllers.ReactivateUser(c) })
	userRoutes.GET(":id/make_admin", func(c *gin.Context) { controllers.MakeUserAdmin(c) })
	userRoutes.GET(":id/remove_admin", func(c *gin.Context) { controllers.RemoveUserAdmin(c) })

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
	sessionRoutes.GET("/dump", func(c *gin.Context) { controllers.DumpSession(c) })
}

func Test(c *gin.Context) {
	c.JSON(200, "This is a Test")
}
