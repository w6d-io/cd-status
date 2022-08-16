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
	"github.com/tektoncd/cli/pkg/formatted"
	"github.com/w6d-io/ci-status/internal/util"
	"github.com/w6d-io/hook"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/w6d-io/ci-status/internal/config"
	"k8s.io/apimachinery/pkg/types"
)

// TaskRunSupervise watches all pod event created by taskrun
func (t *Tekton) TaskRunSupervise(nn types.NamespacedName) error {
	log := t.Log.WithName("TaskRunSupervise").WithValues("object",
		t.PipelineRun.NamespacedName.String(),
		"name", nn.Name)
	w := t.GetWatch("taskrun", nn.Namespace)
	if w == nil {
		return fmt.Errorf("taskrun %s not found", nn.String())
	}
	timeout := time.NewTimer(config.GetTimeout())
	trp := TaskRunPayload{
		NamespacedName: types.NamespacedName{
			Name:      nn.Namespace,
			Namespace: t.PipelineRun.NamespacedName.Namespace,
		},
	}
	t.PipelineRun.UpdateTaskRuns(trp)
	log.WithValues("timeout", config.GetTimeout()).V(1).Info("timeout set")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		log.Info("stop watcher")
		w.Stop()
	}()
	log.V(1).Info("start loop")
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
				if err := hook.Send(ctx, t.PipelineRun, "timeout"); err != nil {
					log.Error(err, "hook failed")
					return err
				}
				return nil
			}
			tr := e.Object.(*tkn.TaskRun)
			t.UpdatePayloadTask(trp.GetTask(tr))

			if err := hook.Send(ctx, t.PipelineRun, "update"); err != nil {
				log.Error(err, "hook failed")
				return err
			}
			if IsTerminated(tr.Status.Conditions) {
				return nil
			}
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

func (TaskRunPayload) GetTask(taskRun *tkn.TaskRun) Task {
	log := logger.WithName("GetTask")
	status := taskRun.Status
	steps := GetSteps(status)
	st, r := Condition(status.Conditions)
	taskName := taskRun.Labels["tekton.dev/pipelineTask"]
	if taskName == "" {
		taskName = taskRun.Spec.TaskRef.Name
	}
	log.V(1).Info("trace", "taskrun_name", taskRun.Name)
	log.V(1).Info("trace", "task_name", taskName)
	return Task{
		TaskRunName:       taskRun.Name,
		Name:              taskName,
		StartTime:         util.UnixMilli(status.StartTime),
		CompletionTime:    util.UnixMilli(status.CompletionTime),
		Duration:          formatted.Duration(status.StartTime, status.CompletionTime),
		Status:            st,
		StartTimeRaw:      status.StartTime,
		CompletionTimeRaw: status.CompletionTime,
		Steps:             steps,
		Message:           r,
	}
}
