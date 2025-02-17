/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

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
	"context"
	json "encoding/json"
	"fmt"

	camelv1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
	applyconfigurationcamelv1 "github.com/apache/camel-k/v2/pkg/client/camel/applyconfiguration/camel/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBuilds implements BuildInterface
type FakeBuilds struct {
	Fake *FakeCamelV1
	ns   string
}

var buildsResource = schema.GroupVersionResource{Group: "camel.apache.org", Version: "v1", Resource: "builds"}

var buildsKind = schema.GroupVersionKind{Group: "camel.apache.org", Version: "v1", Kind: "Build"}

// Get takes name of the build, and returns the corresponding build object, and an error if there is any.
func (c *FakeBuilds) Get(ctx context.Context, name string, options v1.GetOptions) (result *camelv1.Build, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(buildsResource, c.ns, name), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}

// List takes label and field selectors, and returns the list of Builds that match those selectors.
func (c *FakeBuilds) List(ctx context.Context, opts v1.ListOptions) (result *camelv1.BuildList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(buildsResource, buildsKind, c.ns, opts), &camelv1.BuildList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &camelv1.BuildList{ListMeta: obj.(*camelv1.BuildList).ListMeta}
	for _, item := range obj.(*camelv1.BuildList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested builds.
func (c *FakeBuilds) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(buildsResource, c.ns, opts))

}

// Create takes the representation of a build and creates it.  Returns the server's representation of the build, and an error, if there is any.
func (c *FakeBuilds) Create(ctx context.Context, build *camelv1.Build, opts v1.CreateOptions) (result *camelv1.Build, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(buildsResource, c.ns, build), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}

// Update takes the representation of a build and updates it. Returns the server's representation of the build, and an error, if there is any.
func (c *FakeBuilds) Update(ctx context.Context, build *camelv1.Build, opts v1.UpdateOptions) (result *camelv1.Build, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(buildsResource, c.ns, build), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBuilds) UpdateStatus(ctx context.Context, build *camelv1.Build, opts v1.UpdateOptions) (*camelv1.Build, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(buildsResource, "status", c.ns, build), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}

// Delete takes name of the build and deletes it. Returns an error if one occurs.
func (c *FakeBuilds) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(buildsResource, c.ns, name, opts), &camelv1.Build{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBuilds) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(buildsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &camelv1.BuildList{})
	return err
}

// Patch applies the patch and returns the patched build.
func (c *FakeBuilds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *camelv1.Build, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(buildsResource, c.ns, name, pt, data, subresources...), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied build.
func (c *FakeBuilds) Apply(ctx context.Context, build *applyconfigurationcamelv1.BuildApplyConfiguration, opts v1.ApplyOptions) (result *camelv1.Build, err error) {
	if build == nil {
		return nil, fmt.Errorf("build provided to Apply must not be nil")
	}
	data, err := json.Marshal(build)
	if err != nil {
		return nil, err
	}
	name := build.Name
	if name == nil {
		return nil, fmt.Errorf("build.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(buildsResource, c.ns, *name, types.ApplyPatchType, data), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeBuilds) ApplyStatus(ctx context.Context, build *applyconfigurationcamelv1.BuildApplyConfiguration, opts v1.ApplyOptions) (result *camelv1.Build, err error) {
	if build == nil {
		return nil, fmt.Errorf("build provided to Apply must not be nil")
	}
	data, err := json.Marshal(build)
	if err != nil {
		return nil, err
	}
	name := build.Name
	if name == nil {
		return nil, fmt.Errorf("build.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(buildsResource, c.ns, *name, types.ApplyPatchType, data, "status"), &camelv1.Build{})

	if obj == nil {
		return nil, err
	}
	return obj.(*camelv1.Build), err
}
