package router

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
	accessionsRoutes := router.Group("/accessions")
	accessionsRoutes.GET("new", func(c *gin.Context) { controllers.NewAccession(c) })
	accessionsRoutes.POST("", func(c *gin.Context) { controllers.CreateAccession(c) })
	accessionsRoutes.GET("", func(c *gin.Context) { controllers.GetAccessions(c) })
	accessionsRoutes.GET(":id/show", func(c *gin.Context) { controllers.GetAccession(c) })
	accessionsRoutes.GET(":id/edit", func(c *gin.Context) { controllers.EditAccession(c) })
	accessionsRoutes.POST(":id/update", func(c *gin.Context) { controllers.UpdateAccession(c) })
	accessionsRoutes.GET(":id/delete", func(c *gin.Context) { controllers.DeleteAccession(c) })
	accessionsRoutes.GET(":id/slew", func(c *gin.Context) { controllers.SlewAccession(c) })
	accessionsRoutes.POST("slew", func(c *gin.Context) { controllers.CreateAccessionSlew(c) })

	//Repository Group
	repositoryRoutes := router.Group("/repositories")
	repositoryRoutes.GET("", func(c *gin.Context) { controllers.GetRepositories(c) })
	repositoryRoutes.GET(":id/show", func(c *gin.Context) { controllers.GetRepository(c) })
	repositoryRoutes.GET("new", func(c *gin.Context) { controllers.NewRepository(c) })
	repositoryRoutes.POST("", func(c *gin.Context) { controllers.CreateRepository(c) })
	repositoryRoutes.GET(":id/edit", func(c *gin.Context) { controllers.EditRepository(c) })
	repositoryRoutes.POST(":id/update", func(c *gin.Context) { controllers.UpdateRepository(c) })
	repositoryRoutes.GET(":id/delete", func(c *gin.Context) { controllers.DeleteRepository(c) })

	//Resources Group
	resourceRoutes := router.Group("/resources")
	resourceRoutes.GET("", func(c *gin.Context) { controllers.GetResources(c) })
	resourceRoutes.GET(":id/show", func(c *gin.Context) { controllers.GetResource(c) })
	resourceRoutes.GET("new", func(c *gin.Context) { controllers.NewResource(c) })
	resourceRoutes.POST("", func(c *gin.Context) { controllers.CreateResource(c) })
	resourceRoutes.GET(":id/edit", func(c *gin.Context) { controllers.EditResource(c) })
	resourceRoutes.POST(":id/update", func(c *gin.Context) { controllers.UpdateResource(c) })
	resourceRoutes.GET(":id/delete", func(c *gin.Context) { controllers.DeleteResource(c) })

	//Entries Group
	entryRoutes := router.Group("/entries")
	entryRoutes.GET("", func(c *gin.Context) { controllers.GetEntries(c) })
	entryRoutes.GET("new", func(c *gin.Context) { controllers.NewEntry(c) })
	entryRoutes.POST("", func(c *gin.Context) { controllers.CreateEntry(c) })
	entryRoutes.GET(":id/edit", func(c *gin.Context) { controllers.EditEntry(c) })
	entryRoutes.POST(":id/update", func(c *gin.Context) { controllers.UpdateEntry(c) })
	entryRoutes.GET(":id/delete", func(c *gin.Context) { controllers.DeleteEntry(c) })
	entryRoutes.GET(":id/show", func(c *gin.Context) { controllers.GetEntry(c) })
	entryRoutes.GET(":id/previous", func(c *gin.Context) { controllers.GetPreviousEntry(c) })
	entryRoutes.GET(":id/next", func(c *gin.Context) { controllers.GetNextEntry(c) })
	entryRoutes.GET(":id/clone", func(c *gin.Context) { controllers.CloneEntry(c) })
	entryRoutes.POST("find", func(c *gin.Context) { controllers.FindEntry(c) })

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
	userRoutes.GET(":id/show", func(c *gin.Context) { controllers.GetUser(c) })
	userRoutes.GET(":id/edit", func(c *gin.Context) { controllers.EditUser(c) })
	userRoutes.POST("update", func(c *gin.Context) { controllers.UpdateUser(c) })
	userRoutes.GET(":id/allow_api", func(c *gin.Context) { controllers.AllowAPI(c) })
	userRoutes.GET(":id/revoke_api", func(c *gin.Context) { controllers.RevokeAPI(c) })

	//Report Group
	reportRoutes := router.Group("/reports")
	reportRoutes.GET("", func(c *gin.Context) { controllers.ReportsIndex(c) })
	reportRoutes.POST("/range", func(c *gin.Context) { controllers.ReportRange(c) })

	//Search Group
	searchRoutes := router.Group("/search")
	searchRoutes.POST("", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, "Not Implemented") })

	//Session Group
	sessionRoutes := router.Group("/sessions")
	//sessionRoutes.GET("/cookies", func(c *gin.Context) { controllers.GetCookies(c) })
	sessionRoutes.GET("/dump", func(c *gin.Context) { controllers.DumpSession(c) })

	//api v0
	apiV0Routes := router.Group("/api/v0")
	apiV0Routes.POST("users/:user/login", func(c *gin.Context) { controllers.APILogin(c) })
	apiV0Routes.DELETE("logout", func(c *gin.Context) { controllers.APILogout(c) })
	apiV0Routes.GET("", func(c *gin.Context) { controllers.GetV0Index(c) })
	apiV0Routes.GET("repositories/:id", func(c *gin.Context) { controllers.GetRepositoryV0(c) })
	apiV0Routes.GET("repositories", func(c *gin.Context) { controllers.GetRepositoriesV0(c) })
	apiV0Routes.POST("repositories", func(c *gin.Context) { controllers.CreateRepositoryV0(c) })
	apiV0Routes.DELETE("repositories/:id", func(c *gin.Context) { controllers.DeleteRepositoryV0(c) })
	apiV0Routes.GET("repositories/:id/entries", func(c *gin.Context) { controllers.GetRepositoryEntriesV0(c) })
	apiV0Routes.GET("repositories/:id/summary", func(c *gin.Context) { controllers.GetRepositorySummaryV0(c) })
	apiV0Routes.POST("resources", func(c *gin.Context) { controllers.CreateResourceV0(c) })
	apiV0Routes.GET("resources", func(c *gin.Context) { controllers.GetResourcesV0(c) })
	apiV0Routes.GET("resources/:id", func(c *gin.Context) { controllers.GetResourceV0(c) })
	apiV0Routes.DELETE("resources/:id", func(c *gin.Context) { controllers.DeleteResourceV0(c) })
	apiV0Routes.GET("resources/:id/entries", func(c *gin.Context) { controllers.GetResourceEntriesV0(c) })
	apiV0Routes.GET("resources/:id/summary", func(c *gin.Context) { controllers.GetResourceSummaryV0(c) })
	apiV0Routes.POST("accessions", func(c *gin.Context) { controllers.CreateAccessionV0(c) })
	apiV0Routes.DELETE("accessions/:id", func(c *gin.Context) { controllers.DeleteAccessionV0(c) })
	apiV0Routes.GET("accessions", func(c *gin.Context) { controllers.GetAccessionsV0(c) })
	apiV0Routes.GET("accessions/:id", func(c *gin.Context) { controllers.GetAccessionV0(c) })
	apiV0Routes.GET("accessions/:id/entries", func(c *gin.Context) { controllers.GetAccessionEntriesV0(c) })
	apiV0Routes.GET("accessions/:id/summary", func(c *gin.Context) { controllers.GetAccessionSummaryV0(c) })
	apiV0Routes.GET("entries", func(c *gin.Context) { controllers.GetEntriesV0(c) })
	apiV0Routes.GET("entries/:id", func(c *gin.Context) { controllers.GetEntryV0(c) })
	apiV0Routes.PATCH("entries/:id/update_location", func(c *gin.Context) { controllers.UpdateEntryLocationV0(c) })

	// general
	router.NoRoute(func(c *gin.Context) { controllers.NoRoute(c) })
	router.NoMethod(func(c *gin.Context) { c.JSON(http.StatusMethodNotAllowed, "NO METHOD") })

	//error group
	errorRoutes := router.Group("errors")
	errorRoutes.GET("/test", func(c *gin.Context) { controllers.TestError(c) })

}

func Test(c *gin.Context) {
	c.JSON(200, "This is a Test")
}
