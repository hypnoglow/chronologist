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
	"os/signal"
	"sync"
	"syscall"

	"github.com/hypnoglow/chronologist/internal/controller"
	"github.com/hypnoglow/chronologist/internal/grafana"
	"github.com/hypnoglow/chronologist/internal/kube"
	"github.com/hypnoglow/chronologist/internal/zaplog"
)

func main() {
	conf, err := ConfigFromEnvironment()
	if err != nil {
		panic("failed to get config from environment: " + err.Error())
	}

	log, err := zaplog.New(conf.LogFormat, conf.LogLevel)
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	_, kubeClient, err := kube.NewConfigAndClient(conf.KubeConfigPath)
	if err != nil {
		panic("failed to create kubernetes client: " + err.Error())
	}

	grafanaClient := grafana.NewClient(conf.GrafanaAddr, conf.GrafanaAPIKey)

	chronicle := grafana.NewChronicle(grafanaClient, log)

	c, err := controller.New(log, kubeClient, chronicle, controller.Options{
		MaxAge:          conf.ReleaseRevisionMaxAge,
		WatchConfigMaps: conf.WatchConfigMaps,
		WatchSecrets:    conf.WatchSecrets,
	})
	if err != nil {
		panic("failed to create controller: " + err.Error())
	}

	stopCh := make(chan struct{})

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		waitForSignal()

		log.Info("Shutting down ...")
		close(stopCh)
	}()

	c.Run(stopCh)
}

func waitForSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	signal.Stop(signals)
	close(signals)
}
