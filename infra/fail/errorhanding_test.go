package fail_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sre-exporter/infra/fail"
	"sre-exporter/testinfra"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestPanicHandling(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to handle panic with error", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) { panic(fmt.Errorf("some error")) })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"` + fail.ErrUnexpected.Error() + `", "message":"some error", "data": null}`))
	})

	t.Run("should be able to handle panic with other object", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) { panic("some error") })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"` + fail.ErrUnexpected.Error() + `", "message":"some error", "data": null}`))
	})

	t.Run("should be able to handle panic with biz error", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) {
			panic(&demoError{Message: "some message in demo error", Data: 1234})
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(444))
		Expect(body).To(MatchJSON(`{"code":"common.demo", "message":"demo error: some message in demo error", "data": 1234}`))
	})

	t.Run("should not be able to handle panic with nil", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) { panic(nil) })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusOK))
		Expect(body).To(Equal(""))
	})
}

func TestGinErrorHandling(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to handle error in gin.Context.Errors", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) {
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error1")})
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error2")})
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"` + fail.ErrUnexpected.Error() + `", "message":"error2", "data": null}`))
	})

	t.Run("should be able to handle panic error first even gin.Context.Errors is not empty", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) {
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error1")})
			panic("panic error")
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"` + fail.ErrUnexpected.Error() + `", "message":"panic error", "data": null}`))
	})

	t.Run("should handle gin.Context.Errors when panic nil", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) {
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error1")})
			panic(nil)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusInternalServerError))
		Expect(body).To(MatchJSON(`{"code":"` + fail.ErrUnexpected.Error() + `", "message":"error1", "data": null}`))
	})
}

func TestSpecifiedErrorHandling(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should handle common.ErrForbidden", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) {
			_ = c.Error(fail.ErrForbidden)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusForbidden))
		Expect(body).To(MatchJSON(`{"code":"security.forbidden", "message":"access forbidden", "data": null}`))
	})

	t.Run("should handle gorm.ErrRecordNotFound", func(t *testing.T) {
		r := gin.Default()
		r.Use(fail.ErrorHandling())

		r.GET("/", func(c *gin.Context) {
			c.Error(gorm.ErrRecordNotFound)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		status, body, _ := testinfra.ExecuteRequest(req, r)
		Expect(status).To(Equal(http.StatusNotFound))
		Expect(body).To(MatchJSON(`{"code":"common.record_not_found", "message":"record not found", "data": null}`))
	})
}

type demoError struct {
	Message string
	Data    interface{}
}

func (e *demoError) Error() string {
	return fmt.Sprintf("demo error: %s", e.Message)
}
func (e *demoError) Respond() *fail.BizErrorDetail {
	return &fail.BizErrorDetail{
		Status: 444, Code: "common.demo",
		Message: e.Error(), Data: e.Data,
	}
}
