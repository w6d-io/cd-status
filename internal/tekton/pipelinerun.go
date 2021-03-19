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
	"sort"
	"time"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/tektoncd/cli/pkg/formatted"
	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/internal/util"
	"github.com/w6d-io/hook"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis/duck/v1beta1"
)

// PipelineRunSupervise watches all pod event created by pipelinerun
func (t *Tekton) PipelineRunSupervise() error {
	log := t.Log.WithName("PipelineRunSupervise").WithValues("object", t.PipelineRun.NamespacedName.String())
	w := t.GetWatch("pipelinerun", t.PipelineRun.NamespacedName.Name)
	if w == nil {
		return fmt.Errorf("failed to get pipelinerun %s watch", t.PipelineRun.NamespacedName.String())
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
				if err := hook.Send(t.PipelineRun, t.Log, "timeout"); err != nil {
					log.Error(err, "hook failed")
					return err
				}
				return nil
			}
			log.V(1).Info("start sub tasks")
			pr := e.Object.(*tkn.PipelineRun)
			t.PipelineRun.ProjectID = t.ProjectID
			t.PipelineRun.PipelineID = t.PipelineID
			t.PipelineRun.StartTime = util.UnixMilli(pr.Status.StartTime)
			t.SupTasks(pr)
			if IsTerminated(pr.Status.Conditions) {
				t.PipelineRun.CompletionTime = util.UnixMilli(pr.Status.CompletionTime)
				if err := hook.Send(t.PipelineRun, t.Log, "end"); err != nil {
					log.Error(err, "hook failed")
					return err
				}
				return nil
			}
		}
	}
}

// SupTasks loops task watch
func (t *Tekton) SupTasks(pr *tkn.PipelineRun) {
	log := t.Log.WithName("SupTasks").WithValues("object", t.PipelineRun.NamespacedName.String())
	defer log.V(1).Info("SupTasks finished")
	t.PipelineRun.SetCondition(pr.Status.Conditions)
	if err := hook.Send(t.PipelineRun, t.Log, "update"); err != nil {
		log.Error(err, "hook failed")
	}
	currentTask := make(map[string]bool)
	//var wg sync.WaitGroup
	for _, task := range t.PipelineRun.GetTasks(pr) {
		if currentTask[task.Name] {
			log.WithValues("task", task.Name).V(1).Info("skipped")
			continue
		}
		rlog := log.V(1).WithValues("name", task.Name)
		rlog.Info("work")
		currentTask[task.Name] = true
		rlog.V(1).Info("Add wait")
		nn := types.NamespacedName{Namespace: task.TaskRunName, Name: task.Name}
		if err := t.TaskRunSupervise(nn); err != nil {
			log.Error(err, "taskrun supervising")
			return
		}

		//wg.Add(1)

		//go func(wg *sync.WaitGroup, nn types.NamespacedName) {
		//	defer delete(currentTask, nn.Name)
		//	defer removeWait(log, wg)
		//	if err := t.TaskRunSupervise(nn); err != nil {
		//		log.Error(err, "taskrun supervising")
		//		return
		//	}
		//}(&wg, types.NamespacedName{Namespace: task.TaskRunName, Name: task.Name})
	}
	//wg.Wait()
}

// func removeWait(logger logr.Logger, wg *sync.WaitGroup) {
// 	logger.V(1).Info("remove wait")
// 	wg.Done()
// }

func (p *PipelineRunPayload) SetCondition(c v1beta1.Conditions) {
	logger.V(1).Info("SetCondition")
	condition, reason := Condition(c)
	p.Status = condition
	p.Message = reason
}

func (p *PipelineRunPayload) UpdateTaskRuns(taskrun TaskRunPayload) {
	updated := false
	for i, tr := range p.TaskRuns {
		if tr.NamespacedName.String() == taskrun.NamespacedName.String() {
			updated = true
			for _, task := range taskrun.Tasks {
				p.TaskRuns[i].UpdateTask(task)
			}
		}
	}
	if !updated {
		p.TaskRuns = append(p.TaskRuns, taskrun)
	}
}

func (t *Tekton) UpdatePayloadTask(task Task) {
	trp := TaskRunPayload{
		NamespacedName: types.NamespacedName{
			Name:      task.TaskRunName,
			Namespace: t.PipelineRun.NamespacedName.Namespace,
		},
		Tasks: []Task{
			{
				Name:              task.Name,
				Status:            task.Status,
				StartTime:         task.StartTime,
				CompletionTime:    task.CompletionTime,
				StartTimeRaw:      task.StartTimeRaw,
				Duration:          task.Duration,
				Message:           task.Message,
				CompletionTimeRaw: task.CompletionTimeRaw,
				Steps:             task.Steps,
			},
		},
	}
	t.PipelineRun.UpdateTaskRuns(trp)
}

func (p *PipelineRunPayload) GetTasks(pr *tkn.PipelineRun) (ts Tasks) {
	log := logger.WithName("GetTask").WithValues("object", p.NamespacedName.String())
	if len(pr.Status.TaskRuns) > 0 {
		for taskRunName, taskrunStatus := range pr.Status.TaskRuns {
			ts = append(ts, p.GetTask(taskRunName, taskrunStatus))
		}
		sort.Sort(ts)
	}
	log.V(1).WithValues("elements", len(ts)).Info("tasks")
	return
}

func (PipelineRunPayload) GetTask(taskRunName string, taskrunStatus *tkn.PipelineRunTaskRunStatus) Task {
	log := logger.WithValues("method", "GetTask")
	status := taskrunStatus.Status
	st, r := Condition(status.Conditions)
	steps := GetSteps(*status)
	log.V(1).Info("trace", "taskrun_name", taskRunName)
	log.V(1).Info("trace", "task_name", taskrunStatus.PipelineTaskName)
	return Task{
		TaskRunName:       taskRunName,
		Name:              taskrunStatus.PipelineTaskName,
		StartTime:         util.UnixMilli(status.StartTime),
		CompletionTime:    util.UnixMilli(status.CompletionTime),
		Duration:          formatted.Duration(status.StartTime, status.CompletionTime),
		Status:            st,
		StartTimeRaw:      status.StartTime,
		CompletionTimeRaw: status.CompletionTime,
		Message:           r,
		Steps:             steps,
	}
}
