package oss

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewOssClient(t *testing.T) {
	Convey("should be able to build oss client", t, func() {
		os.Setenv("OSS_ENDPOINT", "test-endpoint")
		os.Setenv("OSS_ACCESSKEY", "test-accesskey")
		os.Setenv("OSS_SECRET", "test-secret")

		ossCli := NewOssClient()
		So(*ossCli, ShouldResemble, OssClient{Endpoint: "test-endpoint", AccessKey: "test-accesskey", Secret: "test-secret"})
	})
}
