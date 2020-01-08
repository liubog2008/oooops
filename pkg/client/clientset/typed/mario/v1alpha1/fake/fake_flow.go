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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/liubog2008/oooops/pkg/apis/mario/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFlows implements FlowInterface
type FakeFlows struct {
	Fake *FakeMarioV1alpha1
	ns   string
}

var flowsResource = schema.GroupVersionResource{Group: "mario.oooops.com", Version: "v1alpha1", Resource: "flows"}

var flowsKind = schema.GroupVersionKind{Group: "mario.oooops.com", Version: "v1alpha1", Kind: "Flow"}

// Get takes name of the flow, and returns the corresponding flow object, and an error if there is any.
func (c *FakeFlows) Get(name string, options v1.GetOptions) (result *v1alpha1.Flow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(flowsResource, c.ns, name), &v1alpha1.Flow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Flow), err
}

// List takes label and field selectors, and returns the list of Flows that match those selectors.
func (c *FakeFlows) List(opts v1.ListOptions) (result *v1alpha1.FlowList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(flowsResource, flowsKind, c.ns, opts), &v1alpha1.FlowList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FlowList{ListMeta: obj.(*v1alpha1.FlowList).ListMeta}
	for _, item := range obj.(*v1alpha1.FlowList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested flows.
func (c *FakeFlows) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(flowsResource, c.ns, opts))

}

// Create takes the representation of a flow and creates it.  Returns the server's representation of the flow, and an error, if there is any.
func (c *FakeFlows) Create(flow *v1alpha1.Flow) (result *v1alpha1.Flow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(flowsResource, c.ns, flow), &v1alpha1.Flow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Flow), err
}

// Update takes the representation of a flow and updates it. Returns the server's representation of the flow, and an error, if there is any.
func (c *FakeFlows) Update(flow *v1alpha1.Flow) (result *v1alpha1.Flow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(flowsResource, c.ns, flow), &v1alpha1.Flow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Flow), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFlows) UpdateStatus(flow *v1alpha1.Flow) (*v1alpha1.Flow, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(flowsResource, "status", c.ns, flow), &v1alpha1.Flow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Flow), err
}

// Delete takes name of the flow and deletes it. Returns an error if one occurs.
func (c *FakeFlows) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(flowsResource, c.ns, name), &v1alpha1.Flow{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFlows) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(flowsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.FlowList{})
	return err
}

// Patch applies the patch and returns the patched flow.
func (c *FakeFlows) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Flow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(flowsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Flow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Flow), err
}
