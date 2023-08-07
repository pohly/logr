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
	"log/slog"

	"github.com/go-logr/logr"
)

// NewLogr returns a logr.Logger which writes to the slog.Handler.
//
// In the output the logr verbosity level gets negated, so V(4) becomes
// slog.LevelDebug.
func NewLogr(handler slog.Handler) logr.Logger {
	return logr.New(&slogSink{handler: handler})
}

// NewSlog returns a slog.Handler which writes to the same sink as the logr.Logger.
//
// The returned logger writes all records with level >= slog.LevelError as
// error log entries with LogSink.Error, regardless of the verbosity of the
// logr.Logger. The level of all other records gets reduced by the verbosity
// level of the logr.Logger, so a slog.Logger.Info call gets written with
// slog.LevelDebug when using a logr.Logger where verbosity was modified with
// V(4).
func NewSlog(logger logr.Logger) slog.Handler {
	return &slogHandler{sink: logger.GetSink(), level: slog.Level(logger.GetV())}
}
