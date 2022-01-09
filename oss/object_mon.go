package oss

import (
	"regexp"
	"sre-exporter/infra/metric"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var cycleDuration = time.Hour * 12
var reg = regexp.MustCompile(`\d{8}T\d{6}(Z|\+0800)`)
var bucketName = "fundwit-data"
var ComputeMetricsFunc = ComputeMetrics

func RegisterMetricProvider(ossCli *OssClient) {
	metric.RegisterMetricProvider("oss-backup-objects", func() (metric.Metrics, error) {
		return ComputeMetricsFunc(ossCli)
	})
}

func ComputeMetrics(lister ObjectLister) (metric.Metrics, error) {
	pathIdxMap := map[string]int{}

	paths := []string{}
	counts := []int{}
	lastCycleCounts := []int{}

	marker := ""
	for {
		lsRes, err := lister.ListObjects(bucketName, oss.Marker(marker))
		if err != nil {
			return nil, err
		}

		for _, obj := range lsRes.Objects {
			path, ts := parseBackupObject(obj)
			if path == "" {
				continue
			}

			idx, exist := pathIdxMap[path]
			if !exist {
				idx = len(pathIdxMap)
				pathIdxMap[path] = idx

				paths = append(paths, path)
				counts = append(counts, 0)
				lastCycleCounts = append(lastCycleCounts, 0)
			}

			counts[idx] = counts[idx] + 1

			if time.Since(ts) > cycleDuration {
				continue
			}
			lastCycleCounts[idx] = lastCycleCounts[idx] + 1
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	samples := []metric.Sample{}
	for idx, path := range paths {
		s := metric.Sample{
			Labels: map[string]string{
				"position": "oss://" + bucketName + "/" + path,
				"range":    "all",
			},
			Value: counts[idx],
		}
		samples = append(samples, s)

		s = metric.Sample{
			Labels: map[string]string{
				"position": "oss://" + bucketName + "/" + path,
				"range":    "latest-cycle",
			},
			Value: lastCycleCounts[idx],
		}
		samples = append(samples, s)
	}

	m := metric.Metric{Type: metric.MetricTypeCounter, Name: "backup_data_num", Samples: samples}
	return []metric.Metric{m}, nil
}

func parseBackupObject(obj oss.ObjectProperties) (string, time.Time) {
	parts := strings.Split(obj.Key, "/")
	if len(parts) != 2 {
		return "", time.Time{}
	}

	path := parts[0]
	basename := parts[1]

	timeSections := reg.FindAllString(basename, -1)
	if len(timeSections) != 1 {
		return path, obj.LastModified
	}

	timeString := timeSections[0]
	if strings.HasSuffix(timeString, "Z") {
		ts, _ := time.Parse("20060102T150405Z", timeString)
		return path, ts
	} else {
		ts, _ := time.ParseInLocation("20060102T150405+0800", timeString, time.FixedZone("Asia/Shanghai", 8*60*60))
		return path, ts
	}
}
