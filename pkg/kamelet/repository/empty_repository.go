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

package repository

import (
	"context"

	"github.com/apache/camel-k/v2/pkg/apis/camel/v1alpha1"
)

type emptyKameletRepository struct {
}

func newEmptyKameletRepository() KameletRepository {
	return &emptyKameletRepository{}
}

// Enforce type
var _ KameletRepository = &emptyKameletRepository{}

func (e *emptyKameletRepository) List(_ context.Context) ([]string, error) {
	return nil, nil
}

func (e *emptyKameletRepository) Get(_ context.Context, _ string) (*v1alpha1.Kamelet, error) {
	return nil, nil
}

func (c *emptyKameletRepository) String() string {
	return "Empty[]"
}
