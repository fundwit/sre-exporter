package fail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handle(c)
		c.Next()
	}
}

func handle(c *gin.Context) {
	if ret := recover(); ret != nil {
		err, ok := ret.(error)
		if !ok {
			err = fmt.Errorf("%s", ret)
		}
		HandleError(c, err)
	} else {
		if err := c.Errors.Last(); err != nil {
			HandleError(c, err)
		}
	}
}

func HandleError(c *gin.Context, err error) {
	logrus.Error(err)
	debug.Stack()

	genericErr := err
	var ginErr *gin.Error
	if errors.As(err, &ginErr) {
		genericErr = ginErr.Err
	}

	if bizErr, ok := genericErr.(BizError); ok {
		respond := bizErr.Respond()
		c.JSON(respond.Status, &ErrorBody{Code: respond.Code, Message: respond.Message, Data: respond.Data})
		c.Abort()
		return
	}

	// bad request:  io.EOF (no body).
	if errors.Is(genericErr, io.EOF) {
		c.JSON(http.StatusBadRequest, &ErrorBody{Code: "bad_request.body_not_found", Message: "body not found"})
		c.Abort()
		return
	}
	// bad request: json syntax Error
	if syntaxErr, ok := genericErr.(*json.SyntaxError); ok {
		c.JSON(http.StatusBadRequest, &ErrorBody{Code: "bad_request.invalid_body_format", Message: "invalid body format", Data: syntaxErr.Error()})
		c.Abort()
		return
	}
	// validation failed
	if validationErr, ok := genericErr.(validator.ValidationErrors); ok {
		c.JSON(http.StatusBadRequest, &ErrorBody{Code: "bad_request.validation_failed", Message: "validation failed", Data: validationErr.Error()})
		c.Abort()
		return
	}

	if errors.Is(genericErr, ErrUnauthenticated) {
		c.JSON(http.StatusUnauthorized, &ErrorBody{Code: "common.unauthenticated", Message: "unauthenticated"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, ErrForbidden) {
		c.JSON(http.StatusForbidden, &ErrorBody{Code: "security.forbidden", Message: "access forbidden"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, gorm.ErrRecordNotFound) || errors.Is(genericErr, ErrNotFound) {
		c.JSON(http.StatusNotFound, &ErrorBody{Code: "common.record_not_found", Message: "record not found"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, mysql.ErrInvalidConn) {
		c.JSON(http.StatusServiceUnavailable, &ErrorBody{Code: ErrUnexpected.Error(), Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(500, &ErrorBody{Code: ErrUnexpected.Error(), Message: err.Error()})
	c.Abort()
}
