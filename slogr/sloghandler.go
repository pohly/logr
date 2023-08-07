//go:build go1.21
// +build go1.21

/*
Copyright 2023 The logr Authors.

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

package slogr

import (
	"context"
	"log/slog"

	"github.com/go-logr/logr"
)

type slogHandler struct {
	sink        logr.LogSink
	groupPrefix string
	level       slog.Level
}

var _ slog.Handler = &slogHandler{}

// groupSeparator is used to concatenate WithGroup names and attribute keys.
const groupSeparator = "."

func (l *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return l.sink != nil && (level >= slog.LevelError || l.sink.Enabled(l.levelFromSlog(level)))
}

func (l *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	// No need to check for nil sink here because Handle will only be called
	// when Enabled returned true.

	kvList := make([]any, 0, 2*record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "" {
			kvList = append(kvList, l.addGroupPrefix(attr.Key), attr.Value.Resolve().Any())
		}
		return true
	})
	if record.Level >= slog.LevelError {
		l.sink.Error(nil, record.Message, kvList...)
	} else {
		level := l.levelFromSlog(record.Level)
		l.sink.Info(level, record.Message, kvList...)
	}
	return nil
}

func (l *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if l.sink == nil || len(attrs) == 0 {
		return l
	}
	kvList := make([]any, 0, 2*len(attrs))
	for _, attr := range attrs {
		if attr.Key != "" {
			kvList = append(kvList, l.addGroupPrefix(attr.Key), attr.Value.Resolve().Any())
		}
	}
	copy := *l
	copy.sink = l.sink.WithValues(kvList...)
	return &copy
}

func (l *slogHandler) WithGroup(name string) slog.Handler {
	if l.sink == nil {
		return l
	}
	copy := *l
	copy.groupPrefix = copy.addGroupPrefix(name)
	return &copy
}

func (l *slogHandler) addGroupPrefix(name string) string {
	if l.groupPrefix == "" {
		return name
	}
	return l.groupPrefix + groupSeparator + name
}

func (l *slogHandler) levelFromSlog(level slog.Level) int {
	// Adjust by the per-logger verbosity level.
	level -= l.level
	if level >= 0 {
		// logr has no level lower than 0, so we have to truncate.
		// LogSink implementations might not expect a negative value.
		return 0
	}
	return int(-level)
}
