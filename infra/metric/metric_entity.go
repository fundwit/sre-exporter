package metric

import (
	"fmt"
	"sort"
	"strings"
)

type Metric struct {
	Type    string // counter
	Name    string
	Help    string
	Samples []Sample
}

const (
	MetricTypeCounter = "counter"
)

type Sample struct {
	Labels map[string]string
	Value  interface{}
}

type Metrics []Metric

func (m *Metric) String() string {
	var b strings.Builder
	b.WriteString("# TYPE ")
	b.WriteString(m.Name)
	b.WriteString(" ")
	b.WriteString(m.Type)
	b.WriteString("\n")

	if m.Help != "" {
		b.WriteString("# HELP ")
		b.WriteString(m.Name)
		b.WriteString(" ")
		b.WriteString(strings.ReplaceAll(m.Help, "\n", " "))
		b.WriteString("\n")
	}

	for _, s := range m.Samples {
		b.WriteString(m.Name)

		lastIdx := len(s.Labels) - 1
		if lastIdx >= 0 {
			b.WriteString("{")
			curIdx := 0

			keys := []string{}
			for key := range s.Labels {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				value := s.Labels[key]
				b.WriteString(key)
				b.WriteString("=\"")
				b.WriteString(value)
				b.WriteString("\"")

				if curIdx < lastIdx {
					b.WriteString(",")
				}
				curIdx = curIdx + 1
			}

			b.WriteString("}")
		}

		b.WriteString(" ")
		b.WriteString(fmt.Sprint(s.Value))
		b.WriteString("\n")
	}

	return b.String()
}

func (metrics Metrics) String() string {
	var b strings.Builder
	for _, m := range metrics {
		if len(m.Samples) > 0 {
			b.WriteString(m.String())
		}
	}
	return b.String()
}
