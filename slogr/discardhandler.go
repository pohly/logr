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
)

type discardHandler struct{}

var _ slog.Handler = &discardHandler{}

func (d *discardHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (d *discardHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (d *discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return d
}

func (d *discardHandler) WithGroup(name string) slog.Handler {
	return d
}
