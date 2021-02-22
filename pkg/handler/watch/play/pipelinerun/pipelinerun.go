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
	"github.com/avast/retry-go"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-status/internal/tekton"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

const KIND = "pipelinerun"

// Scan stars the scan of pipeline run tekton resource
func Scan(logger logr.Logger, nn types.NamespacedName, projectID int64, pipelineID int64) error {
	log := logger.WithName("Scan").WithValues("kind", "pipelinerun").
		WithValues("name", nn)
	log.V(1).Info("start")
	defer log.V(1).Info("stop")
	if err := retry.Do(func() error {
		t := &tekton.Tekton{
			ProjectID:  projectID,
			PipelineID: pipelineID,
			Log:        logger,
			PipelineRun: tekton.PipelineRunPayload{
				NamespacedName: types.NamespacedName{
					Name:      nn.Name,
					Namespace: nn.Namespace,
				},
			},
		}
		if err := t.PipelineRunSupervise(); err != nil {
			return err
		}
		return nil
	},
		retry.Delay(3*time.Second),
		retry.Attempts(5),
	); err != nil {
		return err
	}
	return nil
}
