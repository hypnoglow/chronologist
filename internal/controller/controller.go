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

	"github.com/pkg/errors"
	"go.uber.org/zap"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/hypnoglow/chronologist/internal/grafana"
	"github.com/hypnoglow/chronologist/internal/helm"
)

const (
	// maxRetries for an attempt to sync a specific config map with an annotation.
	maxRetries = 5

	// configMapResyncPeriod to resync all config maps.
	configMapResyncPeriod = time.Minute * 10
)

// Controller watches config maps that helm creates for each release
// and creates corresponding annotations in grafana.
type Controller struct {
	log        *zap.Logger
	kubernetes kubernetes.Interface
	grafana    grafana.Annotator

	queue    workqueue.RateLimitingInterface
	informer cache.SharedInformer

	maxAge time.Duration
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

	err := c.syncConfigMap(key.(string))
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

// syncConfigMap method contains logic that is responsible for synchronizing
// a specific config map with a relevant annotation.
func (c *Controller) syncConfigMap(key string) error {
	log := c.log.With(zap.String("configmap", key))

	startTime := time.Now()
	log.Sugar().Infof("Started syncing ConfigMap at %v", startTime.Format(time.RFC3339Nano))
	defer func() {
		log.Sugar().Infof("Finished syncing ConfigMap in %v", time.Since(startTime))
	}()

	// config maps are always named after release revision.
	name, revision, err := c.keyToRelease(key)
	if err != nil {
		return err
	}

	log = log.With(
		zap.String("release", name),
		zap.String("revision", revision),
	)

	q := grafana.GetAnnotationsParams{}
	q.ByRelease(name, revision)

	item, exists, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		return errors.Wrap(err, "get from store by key")
	}
	if !exists {
		log.Sugar().Debugf("ConfigMap has been deleted, deleting grafana annotation")

		aa, err := c.grafana.GetAnnotations(context.TODO(), q)
		if err != nil {
			return err
		}

		var errs []error
		for _, a := range aa {
			log.Sugar().Debugf("Delete grafana annotation id=%d", a.ID)
			if err := c.grafana.DeleteAnnotation(context.TODO(), a.ID); err != nil {
				errs = append(errs, err)
			}
		}

		return utilerrors.NewAggregate(errs)
	}

	cm := item.(*core_v1.ConfigMap)

	relAnn, err := helm.AnnotationFromRawRelease(cm.Data["release"])
	if err != nil {
		return errors.Wrap(err, "create annotation from raw helm release data")
	}

	grafanaAnns, err := c.grafana.GetAnnotations(context.TODO(), q)
	if err != nil {
		return errors.Wrap(err, "get annotations from grafana")
	}

	if len(grafanaAnns) > 1 {
		log.Sugar().Warnf("Release revision has %d annotations. Sync logic for this case is not implemented", len(grafanaAnns))
		// TODO: implement sync logic.
		return nil
	}

	if len(grafanaAnns) < 1 {
		log.Debug("Release revision has no annotations, creating a new one")
		err = c.grafana.SaveAnnotation(
			context.TODO(),
			grafana.AnnotationFromChronologistAnnotation(relAnn),
		)
		return errors.Wrap(err, "create annotation in grafana")
	}

	// Here we got len(grafanaAnns) == 1, which means we need to sync changed
	// config map with corresponding annotation if needed.

	log.Debug("Release revision has one annotation, comparing it with ConfigMap")

	ca, err := grafanaAnns[0].ToChronologistAnnotation()
	if err != nil {
		return errors.Wrap(err, "unmarshal grafana annotation")
	}

	relAnn.GrafanaID = ca.GrafanaID
	diffs := relAnn.Differences(ca)
	if len(diffs) == 0 {
		log.Debug("Annotations are equal, sync is not required")
		return nil
	}

	log.Sugar().Debugf("Found differences: %v. Syncing annotation in grafana", diffs)
	err = c.grafana.SaveAnnotation(
		context.TODO(),
		grafana.AnnotationFromChronologistAnnotation(relAnn),
	)
	if err != nil {
		return errors.Wrap(err, "create annotation")
	}
	return nil
}

// keyToRelease returns release name and revision from config map name.
//
// Config Maps in Helm are named in the way like "foo.v2", where "foo" is the
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

func (c *Controller) addConfigMap(obj interface{}) {
	cm := obj.(*core_v1.ConfigMap)

	// We operate on configmaps that are not outdated.
	if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(cm.CreationTimestamp.Time) {
		c.log.Sugar().Debugf("addConfigMap: ConfigMap %s/%s is too old, skip", cm.Name, cm.Namespace)
		return
	}

	c.log.Sugar().Infof("Adding ConfigMap %s/%s", cm.Namespace, cm.Name)
	c.enqueue(cm)
}

func (c *Controller) updateConfigMap(old, new interface{}) {
	cm := new.(*core_v1.ConfigMap)

	// We operate on configmaps that are not outdated.
	if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(cm.CreationTimestamp.Time) {
		c.log.Sugar().Debugf("updateConfigMap: ConfigMap %s/%s is too old, skip", cm.Name, cm.Namespace)
		return
	}

	c.log.Sugar().Infof("Updating ConfigMap %s/%s", cm.Namespace, cm.Name)
	c.enqueue(cm)
}

func (c *Controller) deleteConfigMap(obj interface{}) {
	cm, ok := obj.(*core_v1.ConfigMap)
	if ok {
		// We operate on configmaps that are not outdated.
		if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(cm.CreationTimestamp.Time) {
			c.log.Sugar().Debugf("deleteConfigMap: ConfigMap %s/%s is too old, skip", cm.Name, cm.Namespace)
			return
		}

		c.log.Sugar().Infof("Deleting ConfigMap %s/%s", cm.Namespace, cm.Name)
		c.enqueue(cm)
		return
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("failed to get object from tombstone %#v", obj))
		return
	}
	_, ok = tombstone.Obj.(*core_v1.ConfigMap)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a ConfigMap %#v", obj))
		return
	}
}

func (c *Controller) enqueue(cm *core_v1.ConfigMap) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cm)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("failed to get key for ConfigMap %s/%s: %v", cm.Namespace, cm.Name, err))
		return
	}

	c.queue.Add(key)
}

// New returns a new controller.
func New(log *zap.Logger, kubernetes kubernetes.Interface, grafana grafana.Annotator, maxAge time.Duration) (*Controller, error) {
	c := &Controller{
		log:        log,
		kubernetes: kubernetes,
		grafana:    grafana,
		maxAge:     maxAge,
	}

	// queue to work on config maps.
	c.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// informer watches for config maps with label OWNER=TILLER
	// and invokes handlers that add those config maps to the queue.
	c.informer = cache.NewSharedInformer(
		// TODO: It would be great if we could filter outdated configmaps here, and not
		// in handler funcs. But this seems impossible currently.
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = "OWNER=TILLER"
				return kubernetes.CoreV1().ConfigMaps(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = "OWNER=TILLER"
				return kubernetes.CoreV1().ConfigMaps(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&core_v1.ConfigMap{},
		configMapResyncPeriod,
	)
	c.informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addConfigMap,
			UpdateFunc: c.updateConfigMap,
			DeleteFunc: c.deleteConfigMap,
		},
	)

	// Hacky stuff.
	utilruntime.ErrorHandlers[0] = func(err error) {
		c.log.Sugar().Errorf("Runtime error: %s", err)
	}

	return c, nil
}
