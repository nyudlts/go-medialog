package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("test the index page", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		req, err := http.NewRequestWithContext(c, "GET", "/", nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 401, recorder.Code)
		assert.Equal(t, "text/html; charset=utf-8", recorder.Header().Get("content-type"))
	})

	t.Run("test login to application", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		form := url.Values{}
		form.Set("email", env.TestCreds.Username)
		form.Add("password_1", env.TestCreds.Password)
		req, err := http.NewRequestWithContext(c, "POST", "/users/authenticate", strings.NewReader(form.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusFound, recorder.Code)

		outFile, _ := os.Create("headers.txt")
		defer outFile.Close()
		writer := bufio.NewWriter(outFile)
		writer.WriteString(fmt.Sprintf("%v", recorder.Result().Header))
		writer.Flush()

	})

	/*
		t.Run("test get index authenticated", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.SetCookie("medialog-sessions",
				"MTcyNTk5NzI2MHxOd3dBTkZkRVFUZEhRazFKU1VWVldqVlRSRlpZUkVWSVYxQkRRMUpHTjAxTFdsVk9OME5HTWpWTVdqSkhURnBPUkVsVVNGWmFWVkU9fGYXoOEZhC7V7FyNzW3NZd0xiPTGNjOqHJg2LMMM5iex",
				2592000, "/", "", false, true,
			)

			req, err := http.NewRequestWithContext(c, "GET", "/", nil)
			if err != nil {
				t.Error(err)
			}
			r.ServeHTTP(recorder, req)
			assert.Equal(t, 200, recorder.Code)
			assert.Equal(t, "text/html; charset=utf-8", recorder.Header().Get("content-type"))

		})
	*/

}
