package pipeline

import (
	"context"
	"errors"
	"testing"
)

type counter struct{ n int }

func stage(name string, fn func(*counter) error) Stage[*counter] {
	return Func[*counter]{StageName: name, StageTitle: name, Fn: func(_ context.Context, c *counter) (string, error) {
		return name + " done", fn(c)
	}}
}

func TestRunRecordsTracesAndSkips(t *testing.T) {
	p := New("test",
		stage("a", func(c *counter) error { c.n++; return nil }),
		stage("b", func(c *counter) error { return Skip("no llm") }),
		stage("c", func(c *counter) error { c.n++; return nil }),
	)
	var streamed []string
	traces, err := p.Run(context.Background(), &counter{}, func(tr Trace) { streamed = append(streamed, tr.Name) })
	if err != nil {
		t.Fatal(err)
	}
	if len(traces) != 3 || traces[1].Status != StatusSkipped || traces[1].Summary != "no llm" {
		t.Fatalf("unexpected traces %+v", traces)
	}
	if len(streamed) != 3 {
		t.Fatalf("expected streaming callbacks for every stage")
	}
}

func TestRunStopsOnFailureAndRecoversPanic(t *testing.T) {
	c := &counter{}
	p := New("test",
		stage("a", func(c *counter) error { panic("boom") }),
		stage("b", func(c *counter) error { c.n++; return nil }),
	)
	traces, err := p.Run(context.Background(), c, nil)
	if err == nil || len(traces) != 1 || traces[0].Status != StatusFailed {
		t.Fatalf("expected failure trace, got %+v %v", traces, err)
	}
	if c.n != 0 {
		t.Fatalf("subsequent stages must not run")
	}

	p = New("test", stage("x", func(c *counter) error { return errors.New("bad") }))
	if _, err := p.Run(context.Background(), c, nil); err == nil {
		t.Fatal("expected error")
	}
}
