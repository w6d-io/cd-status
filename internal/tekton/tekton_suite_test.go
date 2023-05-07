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
    tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
    "k8s.io/apimachinery/pkg/runtime"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
    "testing"

    "k8s.io/client-go/rest"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/envtest"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var (
    ctx       context.Context
    cfg       *rest.Config
    k8sClient client.Client
    testEnv   *envtest.Environment
    scheme    = runtime.NewScheme()
)

func TestTekton(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Tekton Suite")
}

var _ = BeforeSuite(func(done Done) {
    By("bootstrapping test environment")
    testEnv = &envtest.Environment{
        ErrorIfCRDPathMissing: false,
    }

    var err error
    cfg, err = testEnv.Start()
    Expect(err).ToNot(HaveOccurred())
    Expect(cfg).ToNot(BeNil())

    utilruntime.Must(tkn.AddToScheme(scheme))
    k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
    Expect(err).ToNot(HaveOccurred())
    Expect(k8sClient).ToNot(BeNil())

    close(done)
}, 60)

var _ = AfterSuite(func() {
    By("tearing down the test environment")
    err := testEnv.Stop()
    Expect(err).ToNot(HaveOccurred())
})
