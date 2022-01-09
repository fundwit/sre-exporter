package metric

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMetricString(t *testing.T) {
	RegisterTestingT(t)

	t.Run("should be able to compute String", func(t *testing.T) {
		metrics := Metrics{
			{
				Name: "foo", Type: MetricTypeCounter, Help: "this is a foo metric",
				Samples: []Sample{
					{Value: 100, Labels: map[string]string{"k2": "a2", "k1": "a1"}},
					{Value: 200, Labels: map[string]string{"k1": "b1"}},
				},
			},
			{
				Name: "bar", Type: MetricTypeCounter,
				Samples: []Sample{{Value: 300}},
			},
		}

		want := "# TYPE foo counter\n" +
			"# HELP foo this is a foo metric\n" +
			"foo{k1=\"a1\",k2=\"a2\"} 100\n" +
			"foo{k1=\"b1\"} 200\n" +
			"# TYPE bar counter\n" +
			"bar 300\n"

		Expect(metrics.String()).To(Equal(want))
	})
}
