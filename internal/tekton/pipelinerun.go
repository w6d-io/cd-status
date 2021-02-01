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
package tekton

import (
	"context"
	"fmt"
	"sync"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/w6d-io/ci-status/internal/config"
)

// Supervise watches all pod event created by pipelinerun
func (t *Tekton) Supervise() error {
	log := logger.WithName("Supervise").WithValues("object", t.Namespaced.String())
	prw := t.GetWatch("pipelinerun")
	if prw == nil {
		return fmt.Errorf("pipelinerun %s not found", t.Namespaced.String())
	}
	timeout := time.NewTimer(time.Duration(config.GetTimeout()))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer prw.Stop()
	for {
		select {
		case <-timeout.C:
			log.Info("timeout")
			return nil
		case <-ctx.Done():
			log.Info("context Done")
			return nil
		case prc := <-prw.ResultChan():
			if prc.Object == nil {
				log.Info("timeout")
				// TODO notify
				return nil
			}
			t.SupTasks(prc.Object.(*tkn.PipelineRun))
		}
	}
}

// SupTask loops task watch
func (t *Tekton) SupTasks(pr *tkn.PipelineRun) {
	log := logger.WithName("SupTasks").WithValues("object", t.Namespaced.String())
	var wg sync.WaitGroup
	// TODO Notify
	currentTask := make(map[string]bool)
	for _, task := range t.GetTask(pr) {
		if currentTask[task.Name] {
			continue
		}
		rlog := log.V(1).WithValues("name", task.Name)
		rlog.Info("work")
		currentTask[task.Name] = true
		go func(wg *sync.WaitGroup, name string) {
			rlog.V(1).Info("Add waite")
			wg.Add(1)

			wg.Done()
		}(&wg, task.Name)
	}
}
