package main

import (
	"bytes"
	"flag"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/router"
	"github.com/stretchr/testify/assert"
)

var (
	r *gin.Engine
)

func TestAPI(t *testing.T) {

	flag.Parse()

	//set the environment variables

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if sqlite {
		env, err := config.GetSQlite(configuration, environment)
		if err != nil {
			panic(err)
		}
		r, err = router.SetupSQRouter(env, gormDebug)
		if err != nil {
			panic(err)
		}
	} else {

		env, err := config.GetEnvironment(configuration, environment)
		if err != nil {
			panic(err)
		}

		r, err = router.SetupRouter(env, gormDebug)
		if err != nil {
			panic(err)
		}
	}
	t.Run("Test sutup router", func(t *testing.T) {
		t.Logf("context: %v", c)
		t.Logf("router: %v", r)
	})

	t.Run("Test GET login route", func(t *testing.T) {
		req, err := http.NewRequestWithContext(c, "GET", "/users/login", nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(w, req)
		t.Logf("%v", w.Result().Cookies())
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("content-type"))
	})

	t.Run("Test POST login route", func(t *testing.T) {
		var b bytes.Buffer
		w2 := multipart.NewWriter(&b)
		w2.WriteField("email", "admin@nyu.edu")
		w2.WriteField("password_1", "test")
		w2.Close()
		reader := bytes.NewReader(b.Bytes())
		req, err := http.NewRequestWithContext(c, "POST", "/users/authenticate", reader)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", w2.FormDataContentType())
		r.ServeHTTP(w, req)
		t.Logf("%v", w.Header())
		assert.Equal(t, http.StatusOK, w.Code)

		body, _ := io.ReadAll(w.Body)
		title := getTitle(string(body))

		t.Logf("%s", title)
	})

}

var titlePtn = regexp.MustCompile("<title.*/title>")

func getTitle(body string) string {
	matches := titlePtn.FindAllString(body, -1)
	return matches[0]

}

/*

	t.Run("Test GET index route", func(t *testing.T) {
		req, err := http.NewRequestWithContext(c, "GET", "/", nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(w, req)
		t.Logf("%v", w.Result().Cookies())
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("content-type"))
	})


	t.Run("test get Index", func(t *testing.T) {
		req, err := http.NewRequestWithContext(c, "GET", "/", nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(w, req)
		t.Logf("%v", w.Result().Cookies())
	})

	var sessionCookie string
	t.Run("Test POST login route", func(t *testing.T) {
		var b bytes.Buffer
		w2 := multipart.NewWriter(&b)
		w2.WriteField("email", "admin@nyu.edu")
		w2.WriteField("password_1", "test")
		w2.Close()
		reader := bytes.NewReader(b.Bytes())
		req, err := http.NewRequestWithContext(c, "POST", "/users/authenticate", reader)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", w2.FormDataContentType())
		r.ServeHTTP(w, req)
		t.Logf("%v", w.Header())
		header := w.Header().Get("Set-Cookie")
		sessionCookie = strings.Split(strings.Split(header, ";")[0], "=")[1]
		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf(sessionCookie) //placeholder

	})

	t.Run("test get Index", func(t *testing.T) {
		req, err := http.NewRequestWithContext(c, "GET", "/foo", nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(w, req)
		t.Logf("%v", w.Result().Cookies())
	})

	/*


		t.Run("test get Index", func(t *testing.T) {
			req, err := http.NewRequestWithContext(c, "GET", "/", nil)
			if err != nil {
				t.Error(err)
			}
			r.ServeHTTP(w, req)
			t.Logf("%v", w.Result().Cookies())
		})


		var sessionCookie string
		t.Run("Test POST login route", func(t *testing.T) {
			var b bytes.Buffer
			w2 := multipart.NewWriter(&b)
			w2.WriteField("email", "admin@nyu.edu")
			w2.WriteField("password_1", "test")
			w2.Close()
			reader := bytes.NewReader(b.Bytes())
			req, err := http.NewRequestWithContext(c, "POST", "/users/authenticate", reader)
			if err != nil {
				t.Error(err)
			}
			req.Header.Set("Content-Type", w2.FormDataContentType())
			r.ServeHTTP(w, req)
			t.Logf("%v", w.Header())
			header := w.Header().Get("Set-Cookie")
			sessionCookie = strings.Split(strings.Split(header, ";")[0], "=")[1]
			assert.Equal(t, http.StatusOK, w.Code)
			t.Logf(sessionCookie) //placeholder
		})


		t.Run("test logout", func(t *testing.T) {
			req, err := http.NewRequestWithContext(c, "GET", "/users/argh", nil)
			if err != nil {
				t.Error(err)
			}
			r.ServeHTTP(w, req)
			t.Logf("%v", w.Result().Cookies())
		})

		/*




			)


				t.Run("Test GET index", func(t *testing.T) {
					req, err := http.NewRequestWithContext(c, "GET", "/", nil)
					if err != nil {
						t.Error(err)
					}
					r.ServeHTTP(w, req)
					assert.Equal(t, 200, w.Code)
					assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("content-type"))
				})
*/
