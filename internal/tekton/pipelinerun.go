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
	"github.com/w6d-io/ci-status/pkg/hook"
	"knative.dev/pkg/apis/duck/v1beta1"
	"sync"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-status/internal/config"
)

// PipelineRunSupervise watches all pod event created by pipelinerun
func (t *Tekton) PipelineRunSupervise() error {
	log := t.Log.WithName("PipelineRunSupervise").WithValues("object", t.PipelineRun.NamespacedName.String())
	w := t.GetWatch("pipelinerun")
	if w == nil {
		return fmt.Errorf("pipelinerun %s not found", t.PipelineRun.NamespacedName.String())
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
				t.PipelineRun.Status = "timeout"
				if err := hook.Send(t.PipelineRun, log); err != nil {
					log.Error(err, "hook failed")
					return err
				}
				return nil
			}
			t.SupTasks(e.Object.(*tkn.PipelineRun))
			return nil
		}
	}
}

// SupTask loops task watch
func (t *Tekton) SupTasks(pr *tkn.PipelineRun) {
	log := t.Log.WithName("SupTasks").WithValues("object", t.PipelineRun.NamespacedName.String())
	t.PipelineRun.SetCondition(pr.Status.Conditions)
	if err := hook.Send(t.PipelineRun, log); err != nil {
		log.Error(err, "hook failed")
	}
	currentTask := make(map[string]bool)
	var wg sync.WaitGroup
	for _, task := range t.GetTask(pr) {
		if currentTask[task.Name] {
			continue
		}
		rlog := log.V(1).WithValues("name", task.Name)
		rlog.Info("work")
		currentTask[task.Name] = true
		rlog.V(1).Info("Add wait")
		wg.Add(1)
		go func(wg *sync.WaitGroup, name string) {
			defer delete(currentTask, name)
			defer removeWait(log, wg)
			if err := t.TaskRunSupervise(); err != nil {
				log.Error(err, "taskrun supervising")
				return
			}
		}(&wg, task.Name)
	}
	wg.Wait()
}

func removeWait(logger logr.Logger, wg *sync.WaitGroup) {
	logger.V(1).Info("remove wait")
	wg.Done()
}

func (p *PipelineRunPayload) SetCondition(c v1beta1.Conditions) {
	condition, reason := GetStatusReason(c)
	p.Status = condition
	p.Message = reason
}