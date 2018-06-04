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

// ToChronologistAnnotation converts grafana annotation to chronologist
// annotation. This function always returns the same chronologist annotation
// for the same grafana annotation.
func (ga Annotation) ToChronologistAnnotation() (chronologist.Annotation, error) {
	ca := chronologist.Annotation{
		GrafanaID: ga.ID,
		Time:      time.Unix(ga.UNIXMillis/1000, 0).UTC(),
	}

	for _, tag := range ga.Tags {
		switch {
		case strings.HasPrefix(tag, "release_type="):
			rt := strings.TrimPrefix(tag, "release_type=")
			switch rt {
			case chronologist.ReleaseTypeRollout.String():
				ca.ReleaseType = chronologist.ReleaseTypeRollout
			case chronologist.ReleaseTypeRollback.String():
				ca.ReleaseType = chronologist.ReleaseTypeRollback
			default:
				return ca, fmt.Errorf("unknown release type: %v", rt)
			}
		case strings.HasPrefix(tag, "release_status="):
			ca.ReleaseStatus = strings.TrimPrefix(tag, "release_status=")
		case strings.HasPrefix(tag, "release_name="):
			ca.ReleaseName = strings.TrimPrefix(tag, "release_name=")
		case strings.HasPrefix(tag, "release_revision="):
			ca.ReleaseRevision = strings.TrimPrefix(tag, "release_revision=")
		case strings.HasPrefix(tag, "release_namespace="):
			ca.ReleaseNamespace = strings.TrimPrefix(tag, "release_namespace=")
		}
	}

	return ca, nil
}

// AnnotationFromChronologistAnnotation makes a grafana annotation from
// chronologist annotation.
func AnnotationFromChronologistAnnotation(ca chronologist.Annotation) Annotation {
	return Annotation{
		ID:         ca.GrafanaID,
		UNIXMillis: ca.Time.Unix() * 1000,
		Tags: []string{
			"event=release",
			"owner=chronologist",
			"release_type=" + ca.ReleaseType.String(),
			"release_status=" + ca.ReleaseStatus,
			"release_name=" + ca.ReleaseName,
			"release_revision=" + ca.ReleaseRevision,
			"release_namespace=" + ca.ReleaseNamespace,
		},
		Text: fmt.Sprintf("%s release %s", strings.Title(ca.ReleaseType.String()), ca.ReleaseName),
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
		"release_name="+name,
		"release_revision="+revision,
	)
}
