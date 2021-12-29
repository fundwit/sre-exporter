package testinfra

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ExecuteRequest(req *http.Request, engine *gin.Engine) (int, string, *http.Response) {
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	resp := w.Result()
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(bodyBytes), resp
}
