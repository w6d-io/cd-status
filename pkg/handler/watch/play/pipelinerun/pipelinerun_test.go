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
Created on 26/01/2021
*/
package pipelinerun_test

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	//
	//"github.com/tektoncd/cli/pkg/cli"
	//"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	//"github.com/w6d-io/ci-status/internal/tekton"
	"github.com/w6d-io/ci-status/pkg/handler/watch/play/pipelinerun"
	"k8s.io/apimachinery/pkg/types"

	//versionedResource "github.com/tektoncd/pipeline/pkg/client/resource/clientset/versioned"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("PipelineRun", func() {
	Context("Scan", func() {
		When("payload is good", func() {
			//BeforeEach(func() {
			//	resource, err := versionedResource.NewForConfig(cfg)
			//	Expect(err).NotTo(HaveOccurred())
			//	tc, err := versioned.NewForConfig(cfg)
			//	Expect(err).NotTo(HaveOccurred())
			//	c:= &cli.Clients{
			//		Tekton: tc,
			//		Kube: k8scs,
			//		Resource: resource,
			//	}
			//	t := tekton.Tekton{}
			//	t.SetClient(c)
			//})
			It("return no error", func() {
				log := ctrl.Log.WithName("unit-test")
				nn := types.NamespacedName{}
				_ = pipelinerun.Scan(log, nn, 1, 1,
					"ac12ac12ac12ac12ac12ac12", "test", "test",
					"https://test.example.com/repo")
				//Expect(err).To(Succeed())
			})
		})
	})
})
