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
	"errors"
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/hook"
	"github.com/w6d-io/x/logx"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/w6d-io/ci-status/internal/config"
)

func (t *Tekton) TaskSupervise(ctx context.Context, task *Task, w watch.Interface) error {
	log := logx.WithName(ctx, "TaskRunSupervise").WithValues("name", task.TaskRunName)
	log.WithValues("timeout", config.GetTimeout()).V(1).Info("timeout set")
	//var cancel context.CancelFunc
	//ctx, cancel = context.WithTimeout(ctx, config.GetTimeout())
	//defer cancel()
	defer func() {
		log.Info("stop watcher")
		w.Stop()
	}()

	for {
		select {
		case e := <-w.ResultChan():
			log.V(2).Info("watch", "e", e)
			if e.Object == nil {
				log.Info("timeout")
				task.Status = "timeout"
				if err := hook.Send(ctx, t.PipelineRun, "timeout"); err != nil {
					log.Error(err, "hook failed")
					return err
				}
				return nil
			}
			tr := e.Object.(*tkn.TaskRun)
			log.V(1).Info("set actions status")
			task.SetActionsStatus(ctx, &tr.Status)

			if err := hook.Send(ctx, t.PipelineRun, "update"); err != nil {
				log.Error(err, "hook failed")
				return err
			}
			if IsTerminated(tr.Status.Conditions) {
				return nil
			}
		case <-ctx.Done():
			log.Info("context Done", "ctx_message", ctx.Err())
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Info("timeout")
				task.Status = "timeout"
				if err := hook.Send(ctx, t.PipelineRun, "timeout"); err != nil {
					log.Error(err, "hook failed")
					return err
				}
			}
			return nil

		}
	}

}

func (t *TaskRunPayload) UpdateTask(newTask Task) {
	updated := false
	for i, task := range t.Tasks {
		if task.Name == newTask.Name {
			updated = true
			t.Tasks[i] = newTask
			break
		}
	}
	if !updated {
		t.Tasks = append(t.Tasks, newTask)
	}
}
