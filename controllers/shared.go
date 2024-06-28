package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkLogin(c *gin.Context) error {
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return fmt.Errorf("error")
	}
	return nil
}

const UNAUTHORIZED = "Please authenticate to access this service"
