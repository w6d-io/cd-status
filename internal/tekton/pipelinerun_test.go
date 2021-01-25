/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 25/01/2021
*/
package tekton_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/w6d-io/ci-status/internal/tekton"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Tekton", func() {
	Context("Supervise", func() {
		When("timeout is reached", func() {
			It("get out the loop", func() {

			})
		})
	})
	Context("GetWatch", func() {
		BeforeEach(func() {

		})
		AfterEach(func() {

		})
		When("resource exists", func() {
			It("get the watch instance", func() {
				t := &tekton.Tekton{
					Namespaced: types.NamespacedName{
						Name: "pipeline-run-test-1-1",
						Namespace: "default",
					},
				}
				Expect(t.GetWatch()).ToNot(BeNil())
			})
		})
	})
})