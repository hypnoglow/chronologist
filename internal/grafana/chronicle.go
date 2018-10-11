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

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hypnoglow/chronologist/internal/chronologist"
	"github.com/hypnoglow/chronologist/internal/problems"
	"github.com/hypnoglow/chronologist/internal/zaplog"
)

// NewChronicle returns a new Grafana chronicle.
func NewChronicle(grafana Annotator, log *zap.Logger) *Chronicle {
	return &Chronicle{
		grafana: grafana,
		log:     log,
	}
}

// A Chronicle registers release events in Grafana.
type Chronicle struct {
	grafana Annotator
	log     *zap.Logger
}

// Register adds the release event to the chronicle, syncing it with a corresponding
// Grafana annotation.
func (c *Chronicle) Register(ctx context.Context, re chronologist.ReleaseEvent) error {
	log := zaplog.Grasp(ctx, c.log)

	q := GetAnnotationsParams{}
	q.ByRelease(re.Name, re.Revision)

	grafanaAnns, err := c.grafana.GetAnnotations(ctx, q)
	if err != nil {
		return errors.Wrap(err, "get annotations from grafana")
	}

	if len(grafanaAnns) > 1 {
		log.Sugar().Warnf("Found %d annotations for the release event. Sync logic for this case is not implemented", len(grafanaAnns))
		// TODO: implement sync logic.
		return nil
	}

	if len(grafanaAnns) < 1 {
		log.Debug("No annotations found for the release event. Creating a new one")
		err = c.grafana.SaveAnnotation(
			ctx,
			AnnotationFromEvent(0, re),
		)
		return errors.Wrap(err, "create annotation in grafana")
	}

	// Here we got len(grafanaAnns) == 1, which means we need to sync changed
	// release event with corresponding annotation if needed.

	log.Debug("Found one Grafana annotation for the release event. Comparing data")

	re2 := grafanaAnns[0].ToReleaseEvent()

	diffs := re.Differences(re2)
	if len(diffs) == 0 {
		log.Debug("Grafana annotation correctly reflects the release event, sync is not required")
		return nil
	}

	log.Sugar().Debugf("Found differences: %v. Syncing annotation in Grafana", diffs)

	err = c.grafana.SaveAnnotation(
		ctx,
		AnnotationFromEvent(grafanaAnns[0].ID, re),
	)
	if err != nil {
		return errors.Wrap(err, "create annotation")
	}
	return nil
}

// Unregister removes the release event from the chronicle, removing a
// corresponding Grafana annotation.
func (c *Chronicle) Unregister(ctx context.Context, name, revision string) error {
	log := zaplog.Grasp(ctx, c.log)

	q := GetAnnotationsParams{}
	q.ByRelease(name, revision)

	log.Sugar().Debugf("Deleting Grafana annotations related to the release event")

	aa, err := c.grafana.GetAnnotations(ctx, q)
	if err != nil {
		return err
	}

	var errs []error
	for _, a := range aa {
		log.Sugar().Debugf("Delete Grafana annotation id=%d", a.ID)
		if err := c.grafana.DeleteAnnotation(ctx, a.ID); err != nil {
			errs = append(errs, err)
		}
	}

	return problems.NewAggregate(errs)
}
