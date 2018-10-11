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

package grafana

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hypnoglow/chronologist/internal/chronologist"
)

// Annotation represents grafana annotation.
type Annotation struct {
	ID         int      `json:"id,omitempty"`
	UNIXMillis int64    `json:"time"`
	Tags       []string `json:"tags"`
	Text       string   `json:"text"`
}

// ToReleaseEvent converts the grafana annotation to a chronologist release event.
// This function always returns the same chronologist release event for the same
// grafana annotation.
func (a Annotation) ToReleaseEvent() chronologist.ReleaseEvent {
	re := chronologist.ReleaseEvent{
		Time: time.Unix(a.UNIXMillis/1000, 0).UTC(),
	}

	for _, tag := range a.Tags {
		switch {
		case strings.HasPrefix(tag, "release_type="):
			rt := strings.TrimPrefix(tag, "release_type=")
			switch rt {
			case chronologist.ReleaseTypeRollout.String():
				re.Type = chronologist.ReleaseTypeRollout
			case chronologist.ReleaseTypeRollback.String():
				re.Type = chronologist.ReleaseTypeRollback
			default:
				re.Type = chronologist.ReleaseTypeUnknown
			}
		case strings.HasPrefix(tag, "release_status="):
			re.Status = strings.TrimPrefix(tag, "release_status=")
		case strings.HasPrefix(tag, "release_name="):
			re.Name = strings.TrimPrefix(tag, "release_name=")
		case strings.HasPrefix(tag, "release_revision="):
			re.Revision = strings.TrimPrefix(tag, "release_revision=")
		case strings.HasPrefix(tag, "release_namespace="):
			re.Namespace = strings.TrimPrefix(tag, "release_namespace=")
		}
	}

	return re
}

// AnnotationFromEvent assembles a grafana annotation from the chronologist
// release event.
func AnnotationFromEvent(id int, re chronologist.ReleaseEvent) Annotation {
	return Annotation{
		ID:         id,
		UNIXMillis: re.Time.Unix() * 1000,
		Tags: []string{
			"event=release",
			"heritage=chronologist",
			"release_type=" + re.Type.String(),
			"release_status=" + re.Status,
			"release_name=" + re.Name,
			"release_revision=" + re.Revision,
			"release_namespace=" + re.Namespace,
		},
		Text: fmt.Sprintf("%s release %s", strings.Title(re.Type.String()), re.Name),
	}
}

// Annotations is a set of grafana annotations.
type Annotations []Annotation

// Annotator can manage annotations.
type Annotator interface {
	// SaveAnnotation saves annotation, either creating or updating it.
	SaveAnnotation(ctx context.Context, annotation Annotation) error

	// GetAnnotations returns annotations using optional query params.
	GetAnnotations(ctx context.Context, in GetAnnotationsParams) (Annotations, error)

	// DeleteAnnotation deletes annotation by its id.
	DeleteAnnotation(ctx context.Context, id int) error
}

// GetAnnotationsParams represent query parameters for GetAnnotations.
type GetAnnotationsParams struct {
	Tags []string
}

// ByRelease modifies the params to add a filter by specific release name
// and revision.
func (p *GetAnnotationsParams) ByRelease(name, revision string) {
	p.Tags = append(p.Tags,
		"heritage=chronologist",
		"release_name="+name,
		"release_revision="+revision,
	)
}
