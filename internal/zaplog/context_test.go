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

package zaplog_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/hypnoglow/chronologist/internal/zaplog"
)

func TestContext(t *testing.T) {
	log, buf := newTestLogger()

	ctx := context.Background()
	ctx = zaplog.WithFields(ctx, zap.String("a", "b"))
	ctx = zaplog.WithFields(ctx, zap.String("foo", "bar"))
	ctx = zaplog.WithFields(ctx, zap.String("foo", "baz")) // should overwrite
	ctx = zaplog.WithFields(ctx, zap.String("c", "d"))

	zaplog.Grasp(ctx, log).Info("test message")

	var rec testRecord
	err := json.Unmarshal(buf.Bytes(), &rec)
	assert.NoError(t, err)

	assert.Equal(t, rec.A, "b")
	assert.Equal(t, rec.Foo, "baz")
	assert.Equal(t, rec.C, "d")
}

type testRecord struct {
	A       string `json:"a"`
	C       string `json:"c"`
	Foo     string `json:"foo"`
	Level   string `json:"level"`
	Message string `json:"msg"`
}

func newTestLogger() (*zap.Logger, *zaptest.Buffer) {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	buf := &zaptest.Buffer{}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg), buf, zap.DebugLevel,
	)
	return zap.New(core), buf
}
