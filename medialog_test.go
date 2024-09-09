package main

import (
	"net/http"
	"net/http/httptest"
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

}
