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

}
