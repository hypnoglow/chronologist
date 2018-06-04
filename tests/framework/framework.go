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

// Package framework provides conveniences for testing.
package framework

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/services"

	"github.com/hypnoglow/chronologist/internal/grafana"
	"github.com/hypnoglow/chronologist/tests/framework/portforwarder"
)

const (
	tillerNamespace = "kube-system"
	tillerPort      = 44134
	localPort       = 44134
)

type Framework struct {
	mx sync.Mutex

	config    *Config
	k8sClient *kubernetes.Clientset
	k8sConfig *rest.Config

	Helm    *helm.Client
	Grafana *grafana.Client
}

func (f *Framework) SetupTillerTunnel() error {
	f.mx.Lock()
	defer f.mx.Unlock()

	selector := labels.Set{"app": "helm", "name": "tiller"}.AsSelector()
	tunnel, err := portforwarder.New(f.k8sClient, f.k8sConfig, tillerNamespace, selector, tillerPort, localPort)
	if err != nil {
		return err
	}

	f.Helm.Option(helm.Host(fmt.Sprintf("127.0.0.1:%d", tunnel.Local)))
	return nil
}

func (f *Framework) HelmReleaseInstaller(chart, namespace, name string) func() (*services.InstallReleaseResponse, error) {
	// it's weird but otherwise Helm refuses to install the chart.
	b, err := ioutil.ReadFile(filepath.Join(chart, "values.yaml"))
	if err != nil {
		panic(err)
	}

	opts := []helm.InstallOption{
		helm.ReleaseName(name),
		helm.ValueOverrides(b),
	}

	return func() (*services.InstallReleaseResponse, error) {
		return f.Helm.InstallRelease(chart, namespace, opts...)
	}
}

func New() (*Framework, error) {
	f := &Framework{config: &Config{}}
	if err := envconfig.Process("", f.config); err != nil {
		return nil, err
	}

	var err error
	f.k8sConfig, err = clientcmd.BuildConfigFromFlags("", f.config.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	f.k8sClient, err = kubernetes.NewForConfig(f.k8sConfig)
	if err != nil {
		return nil, err
	}

	f.Helm = helm.NewClient(helm.ConnectTimeout(10))
	f.Grafana = grafana.NewClient(f.config.GrafanaAddr, f.config.GrafanaAPIKey)

	return f, nil
}

type Config struct {
	KubeConfigPath string `envconfig:"KUBECONFIG" required:"false"`
	GrafanaAddr    string `envconfig:"GRAFANA_ADDR" required:"true"`
	GrafanaAPIKey  string `envconfig:"GRAFANA_API_KEY" required:"true"`
}
