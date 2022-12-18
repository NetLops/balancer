/*
Copyright 2022 netlops.

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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	balancerv1alpha1 "github.com/netlops/balancer/api/v1alpha1"
	versioned "github.com/netlops/balancer/pkg/client/clientset/versioned"
	internalinterfaces "github.com/netlops/balancer/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/netlops/balancer/pkg/client/listers/balancer/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BalancerInformer provides access to a shared informer and lister for
// Balancers.
type BalancerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BalancerLister
}

type balancerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewBalancerInformer constructs a new informer for Balancer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBalancerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBalancerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredBalancerInformer constructs a new informer for Balancer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBalancerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BalancerV1alpha1().Balancers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BalancerV1alpha1().Balancers(namespace).Watch(context.TODO(), options)
			},
		},
		&balancerv1alpha1.Balancer{},
		resyncPeriod,
		indexers,
	)
}

func (f *balancerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBalancerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *balancerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&balancerv1alpha1.Balancer{}, f.defaultInformer)
}

func (f *balancerInformer) Lister() v1alpha1.BalancerLister {
	return v1alpha1.NewBalancerLister(f.Informer().GetIndexer())
}
