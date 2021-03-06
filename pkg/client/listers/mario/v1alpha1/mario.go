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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MarioLister helps list Marios.
type MarioLister interface {
	// List lists all Marios in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Mario, err error)
	// Marios returns an object that can list and get Marios.
	Marios(namespace string) MarioNamespaceLister
	MarioListerExpansion
}

// marioLister implements the MarioLister interface.
type marioLister struct {
	indexer cache.Indexer
}

// NewMarioLister returns a new MarioLister.
func NewMarioLister(indexer cache.Indexer) MarioLister {
	return &marioLister{indexer: indexer}
}

// List lists all Marios in the indexer.
func (s *marioLister) List(selector labels.Selector) (ret []*v1alpha1.Mario, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Mario))
	})
	return ret, err
}

// Marios returns an object that can list and get Marios.
func (s *marioLister) Marios(namespace string) MarioNamespaceLister {
	return marioNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MarioNamespaceLister helps list and get Marios.
type MarioNamespaceLister interface {
	// List lists all Marios in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Mario, err error)
	// Get retrieves the Mario from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Mario, error)
	MarioNamespaceListerExpansion
}

// marioNamespaceLister implements the MarioNamespaceLister
// interface.
type marioNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Marios in the indexer for a given namespace.
func (s marioNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Mario, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Mario))
	})
	return ret, err
}

// Get retrieves the Mario from the indexer for a given namespace and name.
func (s marioNamespaceLister) Get(name string) (*v1alpha1.Mario, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("mario"), name)
	}
	return obj.(*v1alpha1.Mario), nil
}
