package heldiamgo

import (
	"context"

	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

//开启追踪
var Tracing = true

const TracingChild = 1
const TracingFollowsFrom = 2

//追踪上下文
type TracingContext interface {
	//Start 注意:由于该方法可能会进行连级调用所以不允许返回nil
	Start(operationName string) TracingContext

	//Finish 标记追踪结束
	Finish()
}

type TracingSpan interface {
	//SpanContext
	SpanContext() opentracing.SpanContext

	//context.Context
	Context(ctx context.Context) context.Context

	//Span
	Span() opentracing.Span

	//ChildOf
	ChildOf(operationName string) TracingSpan

	//FollowsFrom
	FollowsFrom(operationName string) TracingSpan

	//日志
	Logf(fields ...opentracinglog.Field)

	//Finish
	Finish()

	//Start
	Start(operationName string) TracingSpan

	//InjectHttpHeader
	InjectHttpHeader(header http.Header)
}

//数据
type TracingSpanImpl struct {
	//operationName   string
	rootSpanContext opentracing.SpanContext
	span            opentracing.Span
}

//SpanContext
func (t *TracingSpanImpl) SpanContext() opentracing.SpanContext {
	return t.span.Context()
}

//context.Context
func (t *TracingSpanImpl) Context(ctx context.Context) context.Context {
	return opentracing.ContextWithSpan(ctx, t.span)
}

//Span
func (t *TracingSpanImpl) Span() opentracing.Span {
	return t.span
}

//ChildOf
func (t *TracingSpanImpl) ChildOf(operationName string) TracingSpan {
	//if t.span != nil {
	//	fmt.Println("childof=>", t.operationName, "=>", operationName)
	//}
	span := opentracing.StartSpan(operationName,
		opentracing.ChildOf(t.SpanContext()))
	return &TracingSpanImpl{
		//operationName: operationName,
		span: span,
	}
}

//FollowsFrom
func (t *TracingSpanImpl) FollowsFrom(operationName string) TracingSpan {
	//if t.span != nil {
	//	fmt.Println("FollowsFrom=>", t.operationName, "=>", operationName)
	//}
	span := opentracing.StartSpan(operationName,
		opentracing.FollowsFrom(t.SpanContext()))
	return &TracingSpanImpl{
		//operationName: operationName,
		span: span,
	}
}

//日志
func (t *TracingSpanImpl) Logf(fields ...opentracinglog.Field) {
	t.span.LogFields(fields...)
}

//Finish
func (t *TracingSpanImpl) Finish() {
	t.span.Finish()
}

//Start
func (t *TracingSpanImpl) Start(operationName string) TracingSpan {
	//if t.span != nil {
	//	fmt.Println("start=>", t.operationName, "=>", operationName)
	//}
	t.span = opentracing.StartSpan(
		operationName,
		ext.RPCServerOption(t.rootSpanContext))
	//t.operationName = operationName
	return t
}

//编码HTTP请求
func (t *TracingSpanImpl) InjectHttpHeader(header http.Header) {
	if header == nil || t.span == nil {
		return
	}
	carrier := opentracing.HTTPHeadersCarrier(header)
	t.span.Tracer().Inject(
		t.span.Context(),
		opentracing.HTTPHeaders,
		carrier)
}

//新建一个追踪
func NewTracingSpanStart(operationName string) TracingSpan {
	//fmt.Println("新数据=>", operationName)
	span := opentracing.StartSpan(operationName)
	return &TracingSpanImpl{
		//operationName: operationName,
		span: span,
	}
}

//解析一个HTTP追踪请求
func NewTracingSpanExtractHttpRequest(req *http.Request) TracingSpan {
	if !Tracing || req == nil {
		return nil
	}
	ret := &TracingSpanImpl{}
	ret.rootSpanContext, _ = opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	ret.Start(req.URL.Path)
	return ret
}
