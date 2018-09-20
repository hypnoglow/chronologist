/*
Copyright 2018 The Chronologist Authors

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

package zaplog

import (
	"context"

	"go.uber.org/zap"
)

type contextFields struct{}

// WithFields returns a context with fields stored into it.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	ff, ok := ctx.Value(contextFields{}).([]zap.Field)
	if !ok {
		return context.WithValue(ctx, contextFields{}, fields)
	}

	// Zap is smart enough to handle duplicate fields for us,
	// so we don't check them ourselves.
	ff = append(ff, fields...)
	return context.WithValue(ctx, contextFields{}, ff)
}

// GetFields returns fields stored in the context.
func GetFields(ctx context.Context) []zap.Field {
	fields, ok := ctx.Value(contextFields{}).([]zap.Field)
	if !ok {
		return nil
	}
	return fields
}

// Grasp makes the logger understand the context. It extracts fields from
// the context, if any, and returns a logger with that fields.
func Grasp(ctx context.Context, log *zap.Logger) *zap.Logger {
	fields, ok := ctx.Value(contextFields{}).([]zap.Field)
	if !ok {
		return log
	}

	return log.With(fields...)
}
