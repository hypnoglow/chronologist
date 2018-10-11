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

func (c *Controller) setupSecretsInformer(kube kubernetes.Interface) {
	// informer watches for secrets with label OWNER=TILLER
	// and invokes handlers that add those secrets to the queue.
	c.informer = cache.NewSharedInformer(
		// TODO: It would be great if we could filter outdated secrets here, and not
		// in handler funcs. But this seems impossible currently.
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = releaseLabelSelector
				return kube.CoreV1().Secrets(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = releaseLabelSelector
				return kube.CoreV1().Secrets(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&core_v1.Secret{},
		releasesResyncPeriod,
	)

	c.informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addSecret,
			UpdateFunc: c.updateSecret,
			DeleteFunc: c.deleteSecret,
		},
	)
}

func (c *Controller) addSecret(obj interface{}) {
	sec := obj.(*core_v1.Secret)

	// We operate on secrets that are not outdated.
	if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(sec.CreationTimestamp.Time) {
		c.log.Sugar().Debugf("addSecret: Secret %s/%s is too old, skip", sec.Name, sec.Namespace)
		return
	}

	c.log.Sugar().Infof("Adding Secret %s/%s", sec.Namespace, sec.Name)
	c.enqueueSecret(sec)
}

func (c *Controller) updateSecret(old, new interface{}) {
	sec := new.(*core_v1.Secret)

	// We operate on secrets that are not outdated.
	if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(sec.CreationTimestamp.Time) {
		c.log.Sugar().Debugf("updateSecret: Secret %s/%s is too old, skip", sec.Name, sec.Namespace)
		return
	}

	c.log.Sugar().Infof("Updating Secret %s/%s", sec.Namespace, sec.Name)
	c.enqueueSecret(sec)
}

func (c *Controller) deleteSecret(obj interface{}) {
	sec, ok := obj.(*core_v1.Secret)
	if ok {
		// We operate on secrets that are not outdated.
		if c.maxAge != 0 && time.Now().Add(-c.maxAge).After(sec.CreationTimestamp.Time) {
			c.log.Sugar().Debugf("deleteSecret: Secret %s/%s is too old, skip", sec.Name, sec.Namespace)
			return
		}

		c.log.Sugar().Infof("Deleting Secret %s/%s", sec.Namespace, sec.Name)
		c.enqueueSecret(sec)
		return
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("failed to get object from tombstone %#v", obj))
		return
	}
	_, ok = tombstone.Obj.(*core_v1.Secret)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Secret %#v", obj))
		return
	}
}

func (c *Controller) enqueueSecret(sec *core_v1.Secret) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(sec)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("failed to get key for Secret %s/%s: %v", sec.Namespace, sec.Name, err))
		return
	}

	c.queue.Add(key)
}

// syncSecret method contains logic that is responsible for synchronizing
// a specific secret with a relevant annotation.
func (c *Controller) syncSecret(key string) error {
	log := c.log.With(zap.String("secret", key))

	startTime := time.Now()
	log.Sugar().Infof("Started syncing Secret at %v", startTime.Format(time.RFC3339Nano))
	defer func() {
		log.Sugar().Infof("Finished syncing Secret in %v", time.Since(startTime))
	}()

	// secrets are always named after release revision.
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

	sec := item.(*core_v1.Secret)

	re, err := helm.EventFromRawRelease(string(sec.Data["release"]))
	if err != nil {
		return errors.Wrap(err, "create a release event from raw helm release data")
	}

	return c.syncReleaseEvent(ctx, re, name, revision)
}
