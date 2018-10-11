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

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/hypnoglow/chronologist/internal/helm"
	"github.com/hypnoglow/chronologist/internal/zaplog"
)

func (c *Controller) setupConfigmapsInformer(kube kubernetes.Interface) {
	// informer watches for configmaps with label OWNER=TILLER
	// and invokes handlers that add those configmaps to the queue.
	c.informer = cache.NewSharedInformer(
		// TODO: It would be great if we could filter outdated configmaps here, and not
		// in handler funcs. But this seems impossible currently.
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = releaseLabelSelector
				return kube.CoreV1().ConfigMaps(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = releaseLabelSelector
				return kube.CoreV1().ConfigMaps(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&core_v1.ConfigMap{},
		releasesResyncPeriod,
	)

	c.informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addConfigMap,
			UpdateFunc: c.updateConfigMap,
			DeleteFunc: c.deleteConfigMap,
		},
	)
}

func (c *Controller) addConfigMap(obj interface{}) {
	cm := obj.(*core_v1.ConfigMap)

	// We operate on configmaps that are not outdated.
	if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(cm.CreationTimestamp.Time) {
		c.log.Sugar().Debugf("addConfigMap: ConfigMap %s/%s is too old, skip", cm.Name, cm.Namespace)
		return
	}

	c.log.Sugar().Infof("Adding ConfigMap %s/%s", cm.Namespace, cm.Name)
	c.enqueueConfigMap(cm)
}

func (c *Controller) updateConfigMap(old, new interface{}) {
	cm := new.(*core_v1.ConfigMap)

	// We operate on configmaps that are not outdated.
	if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(cm.CreationTimestamp.Time) {
		c.log.Sugar().Debugf("updateConfigMap: ConfigMap %s/%s is too old, skip", cm.Name, cm.Namespace)
		return
	}

	c.log.Sugar().Infof("Updating ConfigMap %s/%s", cm.Namespace, cm.Name)
	c.enqueueConfigMap(cm)
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
		c.enqueueConfigMap(cm)
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

func (c *Controller) enqueueConfigMap(cm *core_v1.ConfigMap) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cm)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("failed to get key for ConfigMap %s/%s: %v", cm.Namespace, cm.Name, err))
		return
	}

	c.queue.Add(key)
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

	ctx := zaplog.WithFields(
		context.Background(),
		zap.String("release", name),
		zap.String("revision", revision),
	)

	item, exists, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		return errors.Wrap(err, "get from store by key")
	}
	if !exists {
		return c.deleteReleaseEvent(ctx, name, revision)
	}

	cm := item.(*core_v1.ConfigMap)

	re, err := helm.EventFromRawRelease(cm.Data["release"])
	if err != nil {
		return errors.Wrap(err, "create a release event from raw helm release data")
	}

	return c.syncReleaseEvent(ctx, re, name, revision)
}
