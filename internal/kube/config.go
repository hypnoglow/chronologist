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

// Package kube provides conveniences for working with Kubernetes.
package kube

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewConfigAndClient returns new Kubernetes config and client.
// It reads config from fpath. If fpath is empty, it assumes in-cluster config.
func NewConfigAndClient(fpath string) (*rest.Config, *kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if fpath == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", fpath)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "create kubernetes config")
	}

	client, err := kubernetes.NewForConfig(config)
	return config, client, errors.Wrap(err, "create kubernetes client")
}
