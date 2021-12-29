package fail_test

import (
	"net/http"
	"sre-exporter/infra/fail"
	"testing"

	. "github.com/onsi/gomega"
)

func TestErrBadParam(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return detailed message if param is nil", func(t *testing.T) {
		err := fail.ErrBadParam{Param: "id", InvalidValue: "aaa", Cause: fail.ErrForbidden}
		Expect(err.Error()).To(Equal("invalid id 'aaa'"))
	})

	t.Run("should invoke the Error() function of cause property if cause is not nil", func(t *testing.T) {
		err := fail.ErrBadParam{Cause: fail.ErrForbidden}
		Expect(err.Error()).To(Equal("forbidden"))
	})

	t.Run("should return default message if cause is nil", func(t *testing.T) {
		err := fail.ErrBadParam{}
		Expect(err.Error()).To(Equal("bad param"))
	})
}

func TestErrBadParam_Respond(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return response data as expected", func(t *testing.T) {
		err := fail.ErrBadParam{Param: "id", InvalidValue: "aaa", Cause: fail.ErrForbidden}
		Expect(*err.Respond()).To(Equal(fail.BizErrorDetail{
			Status:  http.StatusBadRequest,
			Code:    "common.bad_param",
			Message: "invalid id 'aaa'",
			Data:    nil,
		}))
	})
}

func TestErrBadParam_Unwrap(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should return response data as expected", func(t *testing.T) {
		err := fail.ErrBadParam{Param: "id", InvalidValue: "aaa", Cause: fail.ErrForbidden}
		Expect(err.Unwrap()).To(Equal(fail.ErrForbidden))

		err = fail.ErrBadParam{}
		Expect(err.Unwrap()).To(BeNil())
	})
}
