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

// Package chronologist provides domain types and logic
// for the application.
package chronologist

import (
	"time"

	"github.com/go-test/deep"
)

// ReleaseType is a type of release.
type ReleaseType string

// String returns release type in a string form.
func (t ReleaseType) String() string {
	return string(t)
}

const (
	// ReleaseTypeRollout is a rollout release type.
	ReleaseTypeRollout ReleaseType = "rollout"

	// ReleaseTypeRollback is a rollback release type.
	ReleaseTypeRollback ReleaseType = "rollback"
)

// Annotation represent an abstract annotation that Chronologist manages.
// Annotation can be serialized or unserialized using different sources.
// See "grafana" and "helm" packages for such functional.
type Annotation struct {
	GrafanaID        int
	Time             time.Time
	ReleaseType      ReleaseType
	ReleaseStatus    string
	ReleaseName      string
	ReleaseRevision  string
	ReleaseNamespace string
}

// Differences compares annotations and returns differences.
func (a Annotation) Differences(ann Annotation) []string {
	return deep.Equal(a, ann)
}
