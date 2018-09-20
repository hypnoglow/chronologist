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

// Package controller encapsulates kubernetes controller logic that powers
// Chronologist.
package controller

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/hypnoglow/chronologist/internal/chronologist"
)

const (
	// maxRetries for an attempt to sync a specific configmap (or secret) with an annotation.
	maxRetries = 5

	// releasesResyncPeriod to resync all configmaps (or secrets).
	releasesResyncPeriod = time.Minute * 10

	// releaseLabelSelector for configmaps (or secrets) created by tiller.
	releaseLabelSelector = "OWNER=TILLER"
)

type releaseBackend string

const (
	backendConfigMaps releaseBackend = "configmaps"
	backendSecrets    releaseBackend = "secrets"
)

// Controller watches configmaps (or secrets) that helm creates for each release
// and creates corresponding annotations in grafana.
type Controller struct {
	log        *zap.Logger
	kubernetes kubernetes.Interface

	queue    workqueue.RateLimitingInterface
	informer cache.SharedInformer

	maxAge time.Duration

	backend releaseBackend

	chronicle chronologist.Chronicle
}

// Options represent controller options.
type Options struct {
	MaxAge          time.Duration
	WatchConfigMaps bool
	WatchSecrets    bool
}

// Run starts the controller.
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stopCh
		c.queue.ShutDown()
	}()

	c.log.Info("Starting controller")
	defer c.log.Info("Shutting down controller")

	switch c.backend {
	case backendConfigMaps:
		c.log.Info("Watch mode: ConfigMaps")
	case backendSecrets:
		c.log.Info("Watch mode: Secrets")
	}

	c.log.Debug("Run informer")
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.informer.Run(stopCh)
	}()

	c.log.Debug("Sync informers cache")
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("failed to sync cache for %v", c.informer))
		return
	}

	c.log.Info("Controller synced and ready")

	wait.Until(c.workerLoop, time.Second, stopCh)
	wg.Wait()
}

func (c *Controller) workerLoop() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	c.log.Sugar().Debugf("Got an item from queue: %s", key.(string))

	var err error
	switch c.backend {
	case backendConfigMaps:
		err = c.syncConfigMap(key.(string))
	case backendSecrets:
		err = c.syncSecret(key.(string))
	default:
		utilruntime.HandleError(fmt.Errorf("release backend is not set up on the controller; this is always a programmer's error"))
	}

	if err == nil {
		c.queue.Forget(key)
		return true
	}

	if c.queue.NumRequeues(key) < maxRetries {
		utilruntime.HandleError(fmt.Errorf("error processing %s (will retry): %v", key, err))
		c.queue.AddRateLimited(key)
		return true
	}

	// Too many retries
	utilruntime.HandleError(fmt.Errorf("error processing %s (giving up): %v", key, err))
	c.queue.Forget(key)

	return true
}

func (c *Controller) syncReleaseEvent(ctx context.Context, re chronologist.ReleaseEvent, name, revision string) error {
	return c.chronicle.Register(ctx, re)
}

func (c *Controller) deleteReleaseEvent(ctx context.Context, name, revision string) error {
	return c.chronicle.Unregister(ctx, name, revision)
}

// keyToRelease returns release name and revision from configmap (or secret) name.
//
// ConfigMaps (or Secrets) in Helm are named in the way like "foo.v2", where "foo" is the
// name of release and "2" is release revision.
func (c *Controller) keyToRelease(key string) (name, revision string, err error) {
	keyParts := strings.SplitN(key, "/", 2)
	if len(keyParts) != 2 {
		return "", "", fmt.Errorf("unknown key format")
	}

	releaseParts := strings.SplitN(keyParts[1], ".", 2)
	if len(releaseParts) != 2 {
		return "", "", fmt.Errorf("unknown key format")
	}
	return releaseParts[0], strings.TrimPrefix(releaseParts[1], "v"), nil
}

// New returns a new controller.
func New(log *zap.Logger, kubernetes kubernetes.Interface, chronicle chronologist.Chronicle, opts Options) (*Controller, error) {
	c := &Controller{
		log:        log,
		kubernetes: kubernetes,
		maxAge:     opts.MaxAge,
		chronicle:  chronicle,
	}

	if opts.WatchConfigMaps && opts.WatchSecrets {
		return nil, fmt.Errorf("incorrect configuration: can watch either configmaps or secrets, not both. Note that in the future Chronologist may be able to watch both")
	}
	if !opts.WatchConfigMaps && !opts.WatchSecrets {
		return nil, fmt.Errorf("incorrect configuration: nothing to watch; need to watch either configmaps or secrets")
	}

	// queue to work on configmaps or secrets.
	c.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	if opts.WatchConfigMaps {
		c.backend = backendConfigMaps
		c.setupConfigmapsInformer(kubernetes)
	} else {
		c.backend = backendSecrets
		c.setupSecretsInformer(kubernetes)
	}

	// Hacky stuff.
	utilruntime.ErrorHandlers[0] = func(err error) {
		c.log.Sugar().Errorf("Runtime error: %s", err)
	}

	return c, nil
}
