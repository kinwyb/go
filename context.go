package heldiamgo

import (
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

type IContext interface {
	TracingContext

	Copy() IContext

	SetTracingType(int)

	Logf(fields ...opentracinglog.Field) IContext

	TracingSpan() TracingSpan
}

//上下文
type Context struct {
	tracing     TracingSpan //追踪数据
	tracingType int         //追踪类型[1=ChildOf,2=FollowsFrom]
	//tracingStart bool              //todo:是否可以通过是否开启追踪来复用一个对象..问题点:多个复用同时开启怎么处理
}

//ChildOf
func ContextChild(parentCtx IContext) IContext {
	if parentCtx == nil {
		return NewContext("ContextChild")
	}
	if Tracing && parentCtx.TracingSpan() != nil {
		ret := parentCtx.Copy()
		ret.SetTracingType(TracingChild)
		return ret
	}
	return parentCtx
}

//FollowsFrom
func ContextFollows(parentCtx IContext) IContext {
	if parentCtx == nil {
		return NewContext("ContextFollows")
	}
	if Tracing && parentCtx.TracingSpan() != nil {
		ret := parentCtx.Copy()
		ret.SetTracingType(TracingFollowsFrom)
		return ret
	}
	return parentCtx
}

//Copy
func (t *Context) Copy() IContext {
	return &Context{
		tracing:     t.tracing,
		tracingType: t.tracingType,
	}
}

func (t *Context) SetTracingType(tracingType int) {
	t.tracingType = tracingType
}

//Finish
func (t *Context) Finish() {
	if t.tracing != nil {
		t.tracing.Finish()
	}
}

//日志
func (t *Context) Logf(fields ...opentracinglog.Field) IContext {
	if t.tracing != nil && Tracing {
		t.tracing.Logf(fields...)
	}
	return t
}

//追踪信息获取,可能返回nil
func (t *Context) TracingSpan() TracingSpan {
	return t.tracing
}

//Start
func (t *Context) Start(operationName string) TracingContext {
	//t.tracingStart = true
	if Tracing {
		if t.tracing == nil || t.tracing.Span() == nil { //没有父级span,生成根span
			t.tracing = NewTracingSpanStart(operationName)
		} else { //有父级span的按类型延伸子级,如果类型为空的不处理
			switch t.tracingType {
			case TracingChild:
				t.tracing = t.tracing.ChildOf(operationName)
			case TracingFollowsFrom:
				t.tracing = t.tracing.FollowsFrom(operationName)
			}
		}
	}
	t.tracingType = 0 //清空追踪类型,往后传递没有指定类型时按之前值往下扩展
	return t
}

//初始化上下文
func NewContext(operationName string) *Context {
	ctx := &Context{}
	if operationName != "" && Tracing {
		ctx.tracing = NewTracingSpanStart(operationName)
	}
	return ctx
}

//初始化上下文
func NewContextWithTracing(tracingSpan TracingSpan) *Context {
	if !Tracing {
		tracingSpan = nil
	}
	return &Context{
		tracing: tracingSpan,
	}
}
