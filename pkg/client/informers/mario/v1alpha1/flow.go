/*
Copyright 2020 The oooops Authors.

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
	time "time"

	mariov1alpha1 "github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	clientset "github.com/liubog2008/oooops/pkg/client/clientset"
	internalinterfaces "github.com/liubog2008/oooops/pkg/client/informers/internalinterfaces"
	v1alpha1 "github.com/liubog2008/oooops/pkg/client/listers/mario/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// FlowInformer provides access to a shared informer and lister for
// Flows.
type FlowInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.FlowLister
}

type flowInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFlowInformer constructs a new informer for Flow type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFlowInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFlowInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFlowInformer constructs a new informer for Flow type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFlowInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MarioV1alpha1().Flows(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MarioV1alpha1().Flows(namespace).Watch(options)
			},
		},
		&mariov1alpha1.Flow{},
		resyncPeriod,
		indexers,
	)
}

func (f *flowInformer) defaultInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFlowInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *flowInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&mariov1alpha1.Flow{}, f.defaultInformer)
}

func (f *flowInformer) Lister() v1alpha1.FlowLister {
	return v1alpha1.NewFlowLister(f.Informer().GetIndexer())
}
