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
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/internal/tekton"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Tekton", func() {
	Context("PipelineRunSupervise", func() {
		BeforeEach(func() {

		})
		When("timeout is reached", func() {
			It("get out the loop", func() {
				err := config.SetTimeout(0)
				Expect(err).To(Succeed())
				err = createPipelineRun("pod-test-1-1")
				Expect(err).NotTo(Succeed())
				t := &tekton.Tekton{
					PipelineRun: &tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pod-test-1-1",
							Namespace: "default",
						},
					},
					Log: ctrl.Log.WithName("test"),
				}
				err = t.PipelineRunSupervise()
				Expect(err).NotTo(Succeed())
				Expect(err.Error()).To(Equal("failed to get pipelinerun default/pod-test-1-1 watch"))
				deletePipelineRun("pr-test-1-1")
			})
		})
	})
	Context("GetWatch", func() {
		When("resource does not exist", func() {
			It("returns nil for pods", func() {
				t := &tekton.Tekton{
					PipelineRun: &tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pod-test-1-1",
							Namespace: "default",
						},
					},
				}
				Expect(t.GetWatch("pods", "test")).To(BeNil())
			})
			It("returns nil for pipelineruns", func() {
				t := &tekton.Tekton{
					PipelineRun: &tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pipeline-run-test-1-1",
							Namespace: "default",
						},
					},
				}
				Expect(t.GetWatch("pipelineruns", "test")).To(BeNil())
			})
			It("returns not nil for taskruns", func() {
				t := &tekton.Tekton{
					PipelineRun: &tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "taskrun-test-1-1",
							Namespace: "default",
						},
					},
				}
				Expect(t.GetWatch("taskruns", "test")).To(BeNil())
			})
			It("returns nil for not supported kind", func() {
				t := &tekton.Tekton{
					PipelineRun: &tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "not-supported-test",
							Namespace: "default",
						},
					},
				}
				Expect(t.GetWatch("not-supported", "test")).To(BeNil())
			})
		})
	})
})

func createPipelineRun(name string) error {
	pr := &tkn.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: tkn.PipelineRunSpec{
			PipelineRef: &tkn.PipelineRef{
				Name: "pipeline-test-1-1",
			},
		},
	}
	if err := k8sClient.Create(context.TODO(), pr); err != nil {
		return err
	}
	return nil
}

func deletePipelineRun(name string) error {
	pr := &tkn.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: tkn.PipelineRunSpec{
			PipelineRef: &tkn.PipelineRef{
				Name: "pipeline-test-1-1",
			},
		},
	}
	if err := k8sClient.Delete(context.TODO(), pr); err != nil {
		return err
	}
	return nil
}
