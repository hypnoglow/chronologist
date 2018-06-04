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

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/helm/pkg/helm"

	"github.com/hypnoglow/chronologist/internal/grafana"
	"github.com/hypnoglow/chronologist/tests/framework"
)

const (
	namespace = "testing"
)

func TestCreateAndDeleteAnnotationForHelmRelease(t *testing.T) {
	fw, err := framework.New()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if err := fw.SetupTillerTunnel(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// On Helm release install, Chronologist should create an annotation in Grafana.

	_, err = fw.HelmReleaseInstaller("./testdata/foo", namespace, "foo")()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer func() {
		// This is just a cleanup so we don't care about errors.
		_, _ = fw.Helm.DeleteRelease("foo", helm.DeletePurge(true))
	}()

	err = wait.Poll(time.Second*2, time.Second*30, func() (done bool, err error) {
		aa, err := fw.Grafana.GetAnnotations(context.TODO(), grafana.GetAnnotationsParams{
			Tags: []string{
				"release_name=foo",
				"release_revision=1",
			},
		})
		if err != nil {
			return false, err
		}
		if len(aa) == 0 {
			return false, nil
		}
		if len(aa) == 1 {
			return true, nil
		}
		return false, fmt.Errorf("unexpected number of annotations: %d", len(aa))
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// On Helm release removal, Chronologist should delete an annotation in Grafana.

	_, err = fw.Helm.DeleteRelease("foo", helm.DeletePurge(true))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = wait.Poll(time.Second*2, time.Second*30, func() (done bool, err error) {
		aa, err := fw.Grafana.GetAnnotations(context.TODO(), grafana.GetAnnotationsParams{
			Tags: []string{
				"release_name=foo",
				"release_revision=1",
			},
		})
		if err != nil {
			return false, err
		}

		return len(aa) == 0, nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
