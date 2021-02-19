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
Created on 28/01/2021
*/
package tekton_test

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-status/internal/tekton"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Tekton", func() {
	Context("GetWatch", func() {
		When("resource does not exist", func() {
			It("returns nil for pods", func() {
				t := &tekton.Tekton{
					PipelineRun: tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pod-test-1-1",
							Namespace: "default",
						},
					},
				}
				Expect(t.GetWatch("pods", "test")).To(BeNil())
			})
		})
	})
	Context("GetTask", func() {
		When("resource does not exist", func() {
			It("returns nil for pods", func() {
				t := &tekton.Tekton{
					PipelineRun: tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pod-test-1-1",
							Namespace: "default",
						},
					},
				}
				pr := &tkn.PipelineRun{
					Status: tkn.PipelineRunStatus{
						PipelineRunStatusFields: tkn.PipelineRunStatusFields{
							StartTime: &metav1.Time{
								Time: time.Now().AddDate(0, 0, -1),
							},
							CompletionTime: &metav1.Time{
								Time: time.Now(),
							},
						},
					},
				}
				Expect(t.PipelineRun.GetTasks(pr)).To(BeNil())
			})
		})
	})
	Context("Tasks sort", func() {
		It("Swap", func() {
			ts := tekton.Tasks{
				{
					Name: "task1",
				},
				{
					Name: "task2",
				},
			}
			Expect(ts.Len()).To(Equal(2))
			Expect(ts.Less(0, 1)).To(Equal(false))
			ts.Swap(0, 1)
			ts[1].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, -1),
			}
			ts.Swap(0, 1)
			ts[0].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, 0),
			}
			ts.Swap(0, 1)
			ts[1].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, 0),
			}
			ts[0].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, -1),
			}
			ts.Swap(0, 1)
		})
	})
})
