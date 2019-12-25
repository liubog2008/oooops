/*
Copyright 2019 The oooops Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	scheme "github.com/liubog2008/oooops/pkg/client/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MariosGetter has a method to return a MarioInterface.
// A group's client should implement this interface.
type MariosGetter interface {
	Marios(namespace string) MarioInterface
}

// MarioInterface has methods to work with Mario resources.
type MarioInterface interface {
	Create(*v1alpha1.Mario) (*v1alpha1.Mario, error)
	Update(*v1alpha1.Mario) (*v1alpha1.Mario, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Mario, error)
	List(opts v1.ListOptions) (*v1alpha1.MarioList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Mario, err error)
	MarioExpansion
}

// marios implements MarioInterface
type marios struct {
	client rest.Interface
	ns     string
}

// newMarios returns a Marios
func newMarios(c *MarioV1alpha1Client, namespace string) *marios {
	return &marios{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mario, and returns the corresponding mario object, and an error if there is any.
func (c *marios) Get(name string, options v1.GetOptions) (result *v1alpha1.Mario, err error) {
	result = &v1alpha1.Mario{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("marios").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Marios that match those selectors.
func (c *marios) List(opts v1.ListOptions) (result *v1alpha1.MarioList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MarioList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("marios").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested marios.
func (c *marios) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("marios").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a mario and creates it.  Returns the server's representation of the mario, and an error, if there is any.
func (c *marios) Create(mario *v1alpha1.Mario) (result *v1alpha1.Mario, err error) {
	result = &v1alpha1.Mario{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("marios").
		Body(mario).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mario and updates it. Returns the server's representation of the mario, and an error, if there is any.
func (c *marios) Update(mario *v1alpha1.Mario) (result *v1alpha1.Mario, err error) {
	result = &v1alpha1.Mario{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("marios").
		Name(mario.Name).
		Body(mario).
		Do().
		Into(result)
	return
}

// Delete takes name of the mario and deletes it. Returns an error if one occurs.
func (c *marios) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("marios").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *marios) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("marios").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mario.
func (c *marios) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Mario, err error) {
	result = &v1alpha1.Mario{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("marios").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
