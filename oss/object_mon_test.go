package oss

import (
	"errors"
	"sre-exporter/infra/metric"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterMetricProvider(t *testing.T) {
	Convey("should be able to compute metrics", t, func() {
		So(metric.GlobalMetricsProviders, ShouldBeEmpty)

		var arg ObjectLister
		ComputeMetricsFunc = func(lister ObjectLister) (metric.Metrics, error) {
			arg = lister
			return metric.Metrics{{Name: "test-metric"}}, errors.New("some-error")
		}

		ossCli := NewOssClient()
		RegisterMetricProvider(ossCli)
		So(metric.GlobalMetricsProviders, ShouldHaveLength, 1)

		f := metric.GlobalMetricsProviders["oss-backup-objects"]
		So(f, ShouldNotBeNil)

		// invoke registered func
		metrics, err := f()
		So(arg, ShouldEqual, ossCli)
		So(err, ShouldBeError, "some-error")
		So(metrics, ShouldResemble, metric.Metrics{{Name: "test-metric"}})
	})
}

func TestComputeMetrics(t *testing.T) {
	Convey("should be able to compute metrics", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish() // invokes assertion

		timeInCycleStr := time.Now().Format("20060102T150405")
		lister := NewMockObjectLister(ctrl)
		// cli.EXPECT().ListObjects(gomock.Eq("fundwit-data"), gomock.Any()).Return(100, errors.New("not exist"))
		lister.EXPECT().
			ListObjects(gomock.Eq("fundwit-data"), gomock.Any()).
			Return(&oss.ListObjectsResult{IsTruncated: true, Objects: []oss.ObjectProperties{
				{Key: "app1/app1-" + timeInCycleStr + "+0800.tar.gz"},
				{Key: "app1/app1-20211011-150405.tar.gz"},     // no timestamp in filename
				{Key: "xxx/app1/app1-20211011-150405.tar.gz"}, // ignored
				{Key: "app2/app2-xx.tar.gz", LastModified: time.Date(2020, 1, 1, 2, 2, 0, 0, time.Local)},
			}}, nil)
		lister.EXPECT().
			ListObjects(gomock.Eq("fundwit-data"), gomock.Any()).
			Return(&oss.ListObjectsResult{IsTruncated: false, Objects: []oss.ObjectProperties{
				{Key: "app3/app3-20211011T150405Z.tar.gz"},
				{Key: "app3/app3-" + timeInCycleStr + "+0800.tar.gz"},
				{Key: "app3/app3-xxxx.tar.gz", LastModified: time.Now()},
			}}, nil)

		want := metric.Metrics{
			{Type: "counter", Name: "backup_data_num", Help: "", Samples: []metric.Sample{
				{Labels: map[string]string{"position": "oss://fundwit-data/app1", "range": "all"}, Value: 2},
				{Labels: map[string]string{"position": "oss://fundwit-data/app1", "range": "latest-cycle"}, Value: 1},
				{Labels: map[string]string{"position": "oss://fundwit-data/app2", "range": "all"}, Value: 1},
				{Labels: map[string]string{"position": "oss://fundwit-data/app2", "range": "latest-cycle"}, Value: 0},
				{Labels: map[string]string{"position": "oss://fundwit-data/app3", "range": "all"}, Value: 3},
				{Labels: map[string]string{"position": "oss://fundwit-data/app3", "range": "latest-cycle"}, Value: 2},
			}},
		}

		metrics, err := ComputeMetrics(lister)
		So(err, ShouldBeNil)
		So(metrics, ShouldResemble, want)
	})

	Convey("should be able to popout error on listObjects", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lister := NewMockObjectLister(ctrl)
		lister.EXPECT().
			ListObjects(gomock.Eq("fundwit-data"), gomock.Any()).
			Return(nil, errors.New("error on list objects"))

		metrics, err := ComputeMetrics(lister)
		So(err, ShouldBeError, "error on list objects")
		So(metrics, ShouldBeNil)
	})
}

func TestParseBackupObject(t *testing.T) {
	Convey("should be able to parse oss object", t, func() {
		t0 := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		path, ts := parseBackupObject(oss.ObjectProperties{Key: "foo/bar.tar.gz", LastModified: t0})
		So(path, ShouldEqual, "foo")
		So(ts, ShouldEqual, t0)

		path, ts = parseBackupObject(oss.ObjectProperties{Key: "foo/foo-20210708T091011+0100.123.tar.gz", LastModified: t0})
		So(path, ShouldEqual, "foo")
		So(ts, ShouldEqual, t0)

		path, ts = parseBackupObject(oss.ObjectProperties{Key: "foo/foo-20210708T091011+0800.123.tar.gz", LastModified: t0})
		So(path, ShouldEqual, "foo")
		So(ts.In(time.UTC), ShouldEqual, time.Date(2021, 7, 8, 1, 10, 11, 0, time.UTC))

		path, ts = parseBackupObject(oss.ObjectProperties{Key: "foo/foo-20210708T091011Z.123.tar.gz", LastModified: t0})
		So(path, ShouldEqual, "foo")
		So(ts, ShouldEqual, time.Date(2021, 7, 8, 9, 10, 11, 0, time.UTC))

		// invalid
		path, ts = parseBackupObject(oss.ObjectProperties{Key: "foo/xxx/bar.tar.gz", LastModified: t0})
		So(path, ShouldBeEmpty)
		So(ts, ShouldEqual, time.Time{})

		path, ts = parseBackupObject(oss.ObjectProperties{Key: "bar.tar.gz", LastModified: t0})
		So(path, ShouldBeEmpty)
		So(ts, ShouldEqual, time.Time{})
	})
}
