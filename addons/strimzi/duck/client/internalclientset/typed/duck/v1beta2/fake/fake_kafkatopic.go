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

	v1beta2 "github.com/apache/camel-k/v2/addons/strimzi/duck/v1beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKafkaTopics implements KafkaTopicInterface
type FakeKafkaTopics struct {
	Fake *FakeKafkaV1beta2
	ns   string
}

var kafkatopicsResource = schema.GroupVersionResource{Group: "kafka.strimzi.io", Version: "v1beta2", Resource: "kafkatopics"}

var kafkatopicsKind = schema.GroupVersionKind{Group: "kafka.strimzi.io", Version: "v1beta2", Kind: "KafkaTopic"}

// Get takes name of the kafkaTopic, and returns the corresponding kafkaTopic object, and an error if there is any.
func (c *FakeKafkaTopics) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.KafkaTopic, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(kafkatopicsResource, c.ns, name), &v1beta2.KafkaTopic{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.KafkaTopic), err
}

// List takes label and field selectors, and returns the list of KafkaTopics that match those selectors.
func (c *FakeKafkaTopics) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.KafkaTopicList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(kafkatopicsResource, kafkatopicsKind, c.ns, opts), &v1beta2.KafkaTopicList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta2.KafkaTopicList{ListMeta: obj.(*v1beta2.KafkaTopicList).ListMeta}
	for _, item := range obj.(*v1beta2.KafkaTopicList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kafkaTopics.
func (c *FakeKafkaTopics) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(kafkatopicsResource, c.ns, opts))

}
