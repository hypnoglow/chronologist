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

package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/hypnoglow/chronologist/internal/zaplog"
)

const prefix = "chronologist"

// Config represents application configuration.
type Config struct {
	// KubeConfigPath is an absolute path to the kubeconfig file.
	KubeConfigPath string `envconfig:"KUBECONFIG" required:"false"`

	GrafanaAddr   string `envconfig:"GRAFANA_ADDR" required:"true"`
	GrafanaAPIKey string `envconfig:"GRAFANA_API_KEY" required:"true"`

	ReleaseRevisionMaxAge time.Duration `envconfig:"RELEASE_REVISION_MAX_AGE" default:"24h"`

	LogFormat zaplog.Format `envconfig:"LOG_FORMAT" default:"json"`
	LogLevel  zaplog.Level  `envconfig:"LOG_LEVEL" default:"info"`

	WatchConfigMaps bool `envconfig:"WATCH_CONFIGMAPS" default:"true"`
	WatchSecrets    bool `envconfig:"WATCH_SECRETS" default:"false"`
}

// ConfigFromEnvironment returns specification loaded from environment
// variables.
func ConfigFromEnvironment() (Config, error) {
	// we do not care if there is no .env file.
	_ = godotenv.Overload()

	var s Config
	err := envconfig.Process(prefix, &s)
	if err != nil {
		return s, err
	}

	s.KubeConfigPath = os.ExpandEnv(s.KubeConfigPath)
	return s, nil
}
