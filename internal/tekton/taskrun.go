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
Created on 06/02/2021
*/
package tekton

import (
	"context"
	"fmt"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-status/internal/config"
	"time"
)

// PipelineRunSupervise watches all pod event created by pipelinerun
func (t *Tekton) TaskRunSupervise() error {
	log := t.Log.WithName("TaskRunSupervise").WithValues("object", t.PipelineRun.NamespacedName.String())
	w := t.GetWatch("taskrun")
	if w == nil {
		return fmt.Errorf("taskrun %s not found", t.PipelineRun.NamespacedName.String())
	}
	timeout := time.NewTimer(config.GetTimeout())
	log.WithValues("timeout", config.GetTimeout()).V(1).Info("timeout set")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		log.Info("stop watcher")
		w.Stop()
	}()
	for {
		select {
		case <-timeout.C:
			log.Info("timeout")
			return nil
		case <-ctx.Done():
			log.Info("context Done")
			return nil
		case e := <-w.ResultChan():
			if e.Object == nil {
				log.Info("timeout")
				// TODO notify
				return nil
			}
			// TODO notify
			t.StepSupervise(e.Object.(*v1beta1.TaskRun))
			return nil
		}
	}
}
