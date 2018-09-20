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

package helm

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/hypnoglow/chronologist/internal/chronologist"
)

// EventFromRelease assembles a chronologist release event from the helm release.
// This function always returns the same event for the same release.
func EventFromRelease(rel *release.Release) (chronologist.ReleaseEvent, error) {
	// We are using LastDeployed field because it is relative only to a specific
	// revision and not to the release itself.
	t, err := ptypes.Timestamp(rel.Info.LastDeployed)
	if err != nil {
		return chronologist.ReleaseEvent{}, errors.Wrap(err, "unserialize timestamp from proto")
	}
	t = t.Truncate(time.Second).UTC()

	rt := chronologist.ReleaseTypeRollout
	if strings.Contains(strings.ToLower(rel.Info.Description), "rollback") {
		rt = chronologist.ReleaseTypeRollback
	}

	return chronologist.ReleaseEvent{
		Time:      t,
		Type:      rt,
		Status:    rel.Info.Status.Code.String(),
		Name:      rel.Name,
		Revision:  strconv.Itoa(int(rel.Version)),
		Namespace: rel.Namespace,
	}, nil
}

// EventFromRawRelease assembles a chronologist release event from the raw helm
// release data. This function always returns the same event for the
// same release.
func EventFromRawRelease(data string) (chronologist.ReleaseEvent, error) {
	rel, err := DecodeRelease(data)
	if err != nil {
		return chronologist.ReleaseEvent{}, errors.Wrap(err, "decode raw release data")
	}

	r, err := EventFromRelease(rel)
	return r, errors.Wrap(err, "create chronologist release event from helm release")
}
