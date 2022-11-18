/*
Copyright 2022 The logr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package benchmarks

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
)

const iterationsPerOp = 100

// 1% of the Info calls are invoked, all of those call the LogSink.
func BenchmarkNewContext1Percent(b *testing.B) {
	expectedCalls := int64(iterationsPerOp) / 100 * int64(b.N)
	ctx := setup(b, 1, expectedCalls)

	// Each iteration is expected to do exactly the same thing,
	// in particular do the same number of allocs.
	for i := 0; i < b.N; i++ {
		// Therefore we repeat newContext a certain number of
		// times. Individual repetitions are allowed to sometimes log
		// and sometimes not, but the overall execution is the same for
		// every outer loop iteration.
		for j := 0; j < iterationsPerOp; j++ {
			newContext(ctx, j, 100, 0)
		}
	}
}

// 100% of the Info calls are invoked, none of those call the LogSink.
func BenchmarkNewContext100PercentDisabled(b *testing.B) {
	expectedCalls := int64(0)
	ctx := setup(b, 1, expectedCalls)

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContext(ctx, j, 1, 2)
		}
	}
}

// 100% of the Info calls are invoked, all of those call the LogSink.
func BenchmarkNewContext100Percent(b *testing.B) {
	expectedCalls := int64(b.N) * iterationsPerOp
	ctx := setup(b, 1, expectedCalls)

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterationsPerOp; j++ {
			newContext(ctx, j, 1, 0)
		}
	}
}

type contextKey1 struct{}
type contextKey2 struct{}

func newContext(ctx context.Context, j, mod, v int) {
	// This is the currently recommended way of adding a value to a context
	// and ensuring that all future log calls include it.  Trace IDs might
	// get handled like this.
	ctx = context.WithValue(ctx, contextKey1{}, 1)
	ctx = context.WithValue(ctx, contextKey2{}, 2)
	useContext(ctx, j, mod, v)
}

func useContext(ctx context.Context, j, mod, v int) {
	if j%mod == 0 {
		logger := logr.FromContextOrDiscard(ctx)
		logger.V(v).Info("ping", "string", "hello world", "int", 1, "float", 1.0)
	}
}

const expectedOutput = `{"logger":"","level":0,"msg":"ping","i":1,"j":2,"string":"hello world","int":1,"float":1}`

func setup(tb testing.TB, v int, expectedCalls int64) context.Context {
	var actualCalls int64
	tb.Cleanup(func() {
		if actualCalls != expectedCalls {
			tb.Errorf("expected %d calls to Info, got %d", expectedCalls, actualCalls)
		}
	})
	logger := funcr.NewJSON(func(actualOutput string) {
		if actualOutput != expectedOutput {
			tb.Fatalf("expected %s, got %s", expectedOutput, actualOutput)
		}
		actualCalls++
	}, funcr.Options{
		FromContextKeys: []funcr.ContextKey{
			funcr.ContextKey{Key: contextKey1{}, Name: "i"},
			funcr.ContextKey{Key: contextKey2{}, Name: "j"},
		},
	})
	return logr.NewContext(context.Background(), logger)
}

func TestFromContext(t *testing.T) {
	expectedCalls := int64(iterationsPerOp) / 100
	ctx := setup(t, 1, expectedCalls)
	newContext(ctx, 0, 1, 0)
}
