package metric

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MetricsProvider = func() (Metrics, error)

var GlobalMetricsProviders = map[string]MetricsProvider{}

var (
	PathMetrics = "/metrics"
)

func RegisterMetricProvider(name string, m MetricsProvider) {
	if m != nil {
		GlobalMetricsProviders[name] = m
	}
}

func RegisterMetricsAPI(r *gin.Engine, middleWares ...gin.HandlerFunc) {
	g := r.Group(PathMetrics, middleWares...)
	g.GET("", handleListMetrics)
}

// @ID metrics
// @Success 200 {string} string metrics data of Prometheus
// @Failure default {object} fail.ErrorBody "error"
// @Router /metrics [get]
func handleListMetrics(c *gin.Context) {
	allMetrics := Metrics{}
	for name, p := range GlobalMetricsProviders {
		m, err := p()
		if err != nil {
			logrus.Warnf("mertic provider '%s', err: %v", name, err)
		}
		allMetrics = append(allMetrics, m...)
	}

	c.String(http.StatusOK, allMetrics.String())
}
