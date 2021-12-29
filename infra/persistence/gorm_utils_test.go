package persistence_test

import (
	"context"
	"sre-exporter/testinfra"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

type TestResource struct {
	ID int
}

func TestGormTracing(t *testing.T) {
	RegisterTestingT(t)

	tracer := mocktracer.New()
	opentracing.SetGlobalTracer(tracer)
	var testDatabase *testinfra.TestDatabase

	t.Run("gorm tracing should be ignored when parent span not found", func(t *testing.T) {
		defer testinfra.GormIntegrateTestTeardown(t, testDatabase)
		testinfra.GormIntegrateTestSetup(t, &testDatabase)
		Expect(testDatabase.GormDB.AutoMigrate(&TestResource{})).To(BeNil())

		tracer.Reset()
		spans := tracer.FinishedSpans()
		Expect(len(spans)).To(Equal(0))

		db := testDatabase.GormDB.WithContext(context.Background())
		r := []TestResource{}
		Expect(db.Find(&r).Error).To(BeNil())
		Expect(len(r)).To(BeZero())

		spans = tracer.FinishedSpans()
		Expect(len(spans)).To(Equal(1))
		Expect(spans[0].ParentID).To(BeZero())
	})

	t.Run("gorm tracing should be work with parent span", func(t *testing.T) {
		defer testinfra.GormIntegrateTestTeardown(t, testDatabase)
		testinfra.GormIntegrateTestSetup(t, &testDatabase)
		Expect(testDatabase.GormDB.AutoMigrate(&TestResource{})).To(BeNil())

		tracer.Reset()

		clientSpan := tracer.StartSpan("client")
		// inject span into context
		ctx := opentracing.ContextWithSpan(context.Background(), clientSpan)
		db := testDatabase.GormDB.WithContext(ctx)

		r := []TestResource{}
		Expect(db.Find(&r).Error).To(BeNil())
		Expect(len(r)).To(BeZero())

		clientSpan.Finish()

		spans := tracer.FinishedSpans()
		Expect(len(spans)).To(Equal(2))
		s0 := spans[1]
		Expect(s0.OperationName).To(Equal("client"))
		Expect(s0.ParentID).To(BeZero())
		Expect(s0.SpanContext.SpanID).ToNot(BeZero())
		Expect(s0.SpanContext.TraceID).ToNot(BeZero())
		Expect(s0.SpanContext.Sampled).To(BeTrue())

		s1 := spans[0]
		Expect(s1.OperationName).To(Equal("query"))
		Expect(s1.ParentID).To(Equal(s0.SpanContext.SpanID))
		Expect(s1.SpanContext.SpanID).ToNot(BeZero())
		Expect(s1.SpanContext.TraceID).To(Equal(s1.SpanContext.TraceID))
		Expect(s1.SpanContext.Sampled).To(BeTrue())
	})
}
