// Package pipeline 是 AI 执行内核的 Pipeline Runtime：按顺序执行 Stage，记录每个 Stage 的执行轨迹。
package pipeline

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Stage 状态
const (
	StatusOK      = "ok"
	StatusSkipped = "skipped"
	StatusFailed  = "failed"
)

// Stage 是 Pipeline 中的一个处理阶段，S 为在各阶段之间传递的共享状态。
type Stage[S any] interface {
	Name() string
	Title() string
	Run(ctx context.Context, state S) (summary string, err error)
}

// Trace 单个 Stage 的执行记录。
type Trace struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	Summary    string `json:"summary"`
	DurationMs int64  `json:"durationMs"`
}

// skipError 表示 Stage 主动跳过（例如大模型未配置），不中断 Pipeline。
type skipError struct{ reason string }

func (e *skipError) Error() string { return e.reason }

// Skip 返回一个"跳过"信号，Runtime 会记录原因并继续执行后续 Stage。
func Skip(reason string) error { return &skipError{reason: reason} }

// Pipeline 由有序 Stage 组成。
type Pipeline[S any] struct {
	Name   string
	Stages []Stage[S]
}

// New 创建 Pipeline。
func New[S any](name string, stages ...Stage[S]) *Pipeline[S] {
	return &Pipeline[S]{Name: name, Stages: stages}
}

// Run 顺序执行全部 Stage；onTrace 在每个 Stage 结束后回调，可用于流式推送进度。
// 任一 Stage 失败即中止并返回错误，已执行的轨迹仍会返回。
func (p *Pipeline[S]) Run(ctx context.Context, state S, onTrace func(Trace)) ([]Trace, error) {
	traces := make([]Trace, 0, len(p.Stages))
	for _, st := range p.Stages {
		if err := ctx.Err(); err != nil {
			return traces, err
		}
		start := time.Now()
		summary, err := runSafe(ctx, st, state)
		tr := Trace{Name: st.Name(), Title: st.Title(), Status: StatusOK, Summary: summary, DurationMs: time.Since(start).Milliseconds()}
		var skip *skipError
		switch {
		case errors.As(err, &skip):
			tr.Status, tr.Summary = StatusSkipped, skip.reason
		case err != nil:
			tr.Status, tr.Summary = StatusFailed, err.Error()
		}
		traces = append(traces, tr)
		if onTrace != nil {
			onTrace(tr)
		}
		if tr.Status == StatusFailed {
			return traces, fmt.Errorf("%s 执行失败: %w", st.Title(), err)
		}
	}
	return traces, nil
}

func runSafe[S any](ctx context.Context, st Stage[S], state S) (summary string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return st.Run(ctx, state)
}

// Func 便于用函数快速声明 Stage。
type Func[S any] struct {
	StageName  string
	StageTitle string
	Fn         func(ctx context.Context, state S) (string, error)
}

func (f Func[S]) Name() string  { return f.StageName }
func (f Func[S]) Title() string { return f.StageTitle }
func (f Func[S]) Run(ctx context.Context, state S) (string, error) {
	return f.Fn(ctx, state)
}
