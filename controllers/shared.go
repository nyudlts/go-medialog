package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func checkLogin(c *gin.Context) error {
	if !isLoggedIn(c) {
		throwError(401, "Please authenticate to access this service", c)
		return fmt.Errorf("error")
	}
	return nil
}
