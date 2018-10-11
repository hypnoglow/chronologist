package grafana_test

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/hypnoglow/chronologist/internal/chronologist"
	"github.com/hypnoglow/chronologist/internal/grafana"
	"github.com/hypnoglow/chronologist/internal/grafana/mocks"
)

// Tests that chronicle creates a new annotation for the release.
func TestChronicle_Register_createAnnotation(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	re := chronologist.ReleaseEvent{
		Time:      time.Date(2019, 01, 02, 15, 4, 5, 0, time.UTC),
		Type:      chronologist.ReleaseTypeRollout,
		Status:    "DEPLOYED",
		Name:      "foo",
		Revision:  "1",
		Namespace: "default",
	}

	ann := mocks.NewAnnotatorMock(t)
	ann.GetAnnotationsMock.
		Expect(context.Background(), grafana.GetAnnotationsParams{
			Tags: []string{
				"heritage=chronologist",
				"release_name=foo",
				"release_revision=1",
			},
		}).
		Return(grafana.Annotations{}, nil)
	ann.SaveAnnotationMock.
		Expect(context.Background(), grafana.Annotation{
			ID:         0,
			UNIXMillis: 1546441445000,
			Tags:       []string{"event=release", "heritage=chronologist", "release_type=rollout", "release_status=DEPLOYED", "release_name=foo", "release_revision=1", "release_namespace=default"},
			Text:       "Rollout release foo",
		}).
		Return(nil)

	cr := grafana.NewChronicle(ann, zap.NewNop())

	err := cr.Register(context.Background(), re)
	assert.NoError(t, err)
}

// Test that chronicle skips the release annotation because it already exists
// and correctly reflects the release event.
func TestChronicle_Register_skipAnnotation(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	re := chronologist.ReleaseEvent{
		Time:      time.Date(2019, 01, 02, 15, 4, 5, 0, time.UTC),
		Type:      chronologist.ReleaseTypeRollout,
		Status:    "DEPLOYED",
		Name:      "foo",
		Revision:  "1",
		Namespace: "default",
	}

	ann := mocks.NewAnnotatorMock(t)
	ann.GetAnnotationsMock.
		Expect(context.Background(), grafana.GetAnnotationsParams{
			Tags: []string{
				"heritage=chronologist",
				"release_name=foo",
				"release_revision=1",
			},
		}).
		Return(grafana.Annotations{{
			ID:         123,
			UNIXMillis: 1546441445000,
			Tags:       []string{"event=release", "heritage=chronologist", "release_type=rollout", "release_status=DEPLOYED", "release_name=foo", "release_revision=1", "release_namespace=default"},
			Text:       "Rollout release foo",
		}}, nil)

	cr := grafana.NewChronicle(ann, zap.NewNop())

	err := cr.Register(context.Background(), re)
	assert.NoError(t, err)
}

// Test that chronicle updates the release annotation because it already exists
// but does not correctly reflect the release event.
func TestChronicle_Register_updateAnnotation(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	re := chronologist.ReleaseEvent{
		Time:      time.Date(2019, 01, 02, 15, 4, 5, 0, time.UTC),
		Type:      chronologist.ReleaseTypeRollout,
		Status:    "DEPLOYED",
		Name:      "foo",
		Revision:  "1",
		Namespace: "default",
	}

	ann := mocks.NewAnnotatorMock(t)
	ann.GetAnnotationsMock.
		Expect(context.Background(), grafana.GetAnnotationsParams{
			Tags: []string{
				"heritage=chronologist",
				"release_name=foo",
				"release_revision=1",
			},
		}).
		Return(grafana.Annotations{{
			ID:         123,
			UNIXMillis: 1546441439000,
			Tags:       []string{"event=release", "heritage=chronologist", "release_type=rollout", "release_status=DEPLOYED", "release_name=foo", "release_revision=1", "release_namespace=default"},
			Text:       "Rollout release foo",
		}}, nil)
	ann.SaveAnnotationMock.
		Expect(context.Background(), grafana.Annotation{
			ID:         123,
			UNIXMillis: 1546441445000,
			Tags:       []string{"event=release", "heritage=chronologist", "release_type=rollout", "release_status=DEPLOYED", "release_name=foo", "release_revision=1", "release_namespace=default"},
			Text:       "Rollout release foo",
		}).
		Return(nil)

	cr := grafana.NewChronicle(ann, zap.NewNop())

	err := cr.Register(context.Background(), re)
	assert.NoError(t, err)
}

func TestChronicle_Unregister(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	re := chronologist.ReleaseEvent{
		Time:      time.Date(2019, 01, 02, 15, 4, 5, 0, time.UTC),
		Type:      chronologist.ReleaseTypeRollout,
		Status:    "DEPLOYED",
		Name:      "foo",
		Revision:  "1",
		Namespace: "default",
	}

	ann := mocks.NewAnnotatorMock(t)
	ann.GetAnnotationsMock.
		Expect(context.Background(), grafana.GetAnnotationsParams{
			Tags: []string{
				"heritage=chronologist",
				"release_name=foo",
				"release_revision=1",
			},
		}).
		Return(grafana.Annotations{{
			ID:         123,
			UNIXMillis: 1546441445000,
			Tags:       []string{"event=release", "heritage=chronologist", "release_type=rollout", "release_status=DEPLOYED", "release_name=foo", "release_revision=1", "release_namespace=default"},
			Text:       "Rollout release foo",
		}}, nil)
	ann.DeleteAnnotationMock.
		Expect(context.Background(), 123).
		Return(nil)

	cr := grafana.NewChronicle(ann, zap.NewNop())

	err := cr.Unregister(context.Background(), re.Name, re.Revision)
	assert.NoError(t, err)
}
