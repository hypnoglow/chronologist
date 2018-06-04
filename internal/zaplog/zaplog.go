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

// Package zaplog provides conveniences for zap logger.
package zaplog

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level defines log level.
type Level struct {
	zap.AtomicLevel
}

// Format defines log format.
type Format string

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Format) UnmarshalText(text []byte) error {
	txt := string(text)
	switch txt {
	case "console", "text":
		*f = FormatConsole
	case "json":
		*f = FormatJSON
	default:
		return fmt.Errorf("unknown format: %q", txt)
	}

	return nil
}

// String implement fmt.Stringer.
func (f Format) String() string {
	return string(f)
}

const (
	// FormatConsole is console log format.
	FormatConsole Format = "console"

	// FormatJSON is JSON log format.
	FormatJSON Format = "json"
)

// New returns a new zap logger.
func New(format Format, level Level) (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig = encoderConfig(format)
	logConfig.Encoding = format.String()
	logConfig.Level = level.AtomicLevel
	return logConfig.Build()

}

func encoderConfig(format Format) zapcore.EncoderConfig {
	var le zapcore.LevelEncoder
	switch format {
	case FormatConsole:
		le = zapcore.CapitalColorLevelEncoder
	default:
		le = zapcore.LowercaseLevelEncoder
	}

	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    le,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
