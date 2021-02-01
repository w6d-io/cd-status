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
Created on 24/01/2021
*/
package pipelinerun

import (
	"github.com/w6d-io/ci-status/internal/tekton"
	"github.com/w6d-io/ci-status/pkg/watch"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PipelineRun struct {
}

var (
	pipelineRun watch.Interface
	logger      = ctrl.Log.WithName("api").WithName("watch").WithName("pipelinerun")
)

const KIND = "pipelinerun"

func init() {
	pipelineRun = new(PipelineRun)
	watch.AddWatcher(KIND, pipelineRun.Scan)
}

func (p PipelineRun) Scan(payload watch.Payload) error {
	log := logger.WithName("Scan").
		WithValues("kind", payload.Object.Kind).
		WithValues("name", payload.Object.NamespacedName)
	log.V(1).Info("start")
	nn := payload.Object.NamespacedName
	t := tekton.Tekton{
		ProjectID:  payload.ProjectID,
		PipelineID: payload.PipelineID,
		Namespaced: types.NamespacedName{
			Name:      nn.Name,
			Namespace: nn.Namespace,
		},
	}
	if err := t.Supervise(); err != nil {
		return err
	}
	return nil
}
