package metric

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sre-exporter/infra/fail"
	"sre-exporter/testinfra"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetirsAPI(t *testing.T) {
	Convey("test metrics API", t, func() {
		router := gin.Default()
		router.Use(fail.ErrorHandling())
		RegisterMetricsAPI(router)

		Convey("should be able to return metrics", func() {
			RegisterMetricProvider("p1", func() (Metrics, error) {
				return Metrics{{Name: "p1", Type: MetricTypeCounter, Samples: []Sample{{Value: 100}}}}, nil
			})
			RegisterMetricProvider("p2", func() (Metrics, error) {
				return Metrics{{Name: "p2", Type: MetricTypeCounter, Samples: []Sample{{Value: 200}}}}, nil
			})
			RegisterMetricProvider("p3", func() (Metrics, error) {
				return nil, errors.New("some error")
			})

			want := "# TYPE p1 counter\n" +
				"p1 100\n" +
				"# TYPE p2 counter\n" +
				"p2 200\n"

			req := httptest.NewRequest(http.MethodGet, PathMetrics, nil)
			status, body, _ := testinfra.ExecuteRequest(req, router)
			So(status, ShouldEqual, 200)
			So(body, ShouldEqual, want)
		})
	})
}
