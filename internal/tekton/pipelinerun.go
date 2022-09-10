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
    "github.com/w6d-io/x/logx"
    "time"

    tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

    "github.com/w6d-io/ci-status/internal/config"
    "github.com/w6d-io/ci-status/internal/util"
    "github.com/w6d-io/hook"
    "k8s.io/apimachinery/pkg/types"
    "knative.dev/pkg/apis/duck/v1beta1"
)

// PipelineRunSupervise watches all pod event created by pipelinerun
func (t *Tekton) PipelineRunSupervise(ctx context.Context) error {
    log := logx.WithName(ctx, "PipelineRunSupervise").WithValues("object", t.PipelineRun.NamespacedName().String())
    w := t.GetWatch(ctx, "pipelinerun", t.PipelineRun.NamespacedName().Name)
    if w == nil {
        return fmt.Errorf("failed to get pipelinerun %s watch", t.PipelineRun.NamespacedName().String())
    }
    timeout := time.NewTimer(config.GetTimeout())
    log.WithValues("timeout", config.GetTimeout(), "message", ctx.Err()).V(1).Info("timeout set")

    var cancel context.CancelFunc
    ctx, cancel = context.WithTimeout(ctx, config.GetTimeout())
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
            log.Info("context Done", "message", ctx.Err())
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
            log.V(1).Info("start sub tasks")
            pr := e.Object.(*tkn.PipelineRun)
            t.SupTasks(ctx, pr)
            if IsTerminated(pr.Status.Conditions) {
                t.PipelineRun.EndTime = util.UnixMilli(pr.Status.CompletionTime)
                if err := hook.Send(ctx, t.PipelineRun, "end"); err != nil {
                    log.Error(err, "hook failed")
                    return err
                }
                return nil
            }
        }
    }
}

// SupTasks loops task watch
func (t *Tekton) SupTasks(ctx context.Context, pr *tkn.PipelineRun) {
    log := logx.WithName(ctx, "SupTasks").WithValues("object", t.PipelineRun.NamespacedName().String())
    defer log.V(1).Info("SupTasks finished")
    t.PipelineRun.SetCondition(ctx, pr.Status.Conditions)
    if err := hook.Send(context.Background(), t.PipelineRun, "update"); err != nil {
        log.Error(err, "hook failed")
    }
    currentTask := make(map[string]bool)

    for _, task := range t.PipelineRun.GetTasks(ctx, pr) {
        if currentTask[task.Name] {
            log.WithValues("task", task.Name).V(1).Info("skipped")
            continue
        }
        rlog := log.V(1).WithValues("name", task.Name, "task_id", task.ID)
        rlog.Info("work")
        currentTask[task.Name] = true
        rlog.V(1).Info("Add wait")

        go func(ctx context.Context, task *Task) {
            defer log.V(1).Info("TaskSupervise finished", "taskrun_name", task.TaskRunName)
            w := t.GetWatch(ctx, "taskrun", task.TaskRunName)
            if w != nil {
                if err := t.TaskSupervise(ctx, task, w); err != nil {
                    log.Error(err, "routine failed", "taskrun_name", task.TaskRunName)
                }
            }

        }(ctx, task)
    }
}

func (p *PipelineRunPayload) SetCondition(ctx context.Context, c v1beta1.Conditions) {
    logx.WithName(ctx, "PipelineRunPayload.SetCondition").V(1).Info("SetCondition")
    condition, reason := GetStatusFromCondition(c)
    p.Status = condition
    p.Reason = reason
}

func (p *PipelineRunPayload) GetTasks(ctx context.Context, pr *tkn.PipelineRun) (ts Tasks) {
    log := logx.WithName(ctx, "GetTask").WithValues("object", p.NamespacedName().String())
    if len(pr.Status.TaskRuns) > 0 {
        for taskRunName, taskrunStatus := range pr.Status.TaskRuns {
            if task := p.GetTask(ctx, taskRunName, taskrunStatus); task != nil {
                ts = append(ts, task)
            }
        }
        //sort.Sort(ts)
    }
    log.V(1).WithValues("elements", len(ts)).Info("tasks")
    return
}

func (p *PipelineRunPayload) GetTask(ctx context.Context, taskRunName string, taskrunStatus *tkn.PipelineRunTaskRunStatus) *Task {
    log := logx.WithName(ctx, "tekton.GetTask")
    status := taskrunStatus.Status
    st, _ := GetStatusFromCondition(status.Conditions)

    log.V(1).Info("trace", "task_name", taskrunStatus.PipelineTaskName)
    var task *Task

    for i, stage := range p.Stages {
        for j, t := range stage.Tasks {
            log.V(1).Info("trace", "t.ID", t.ID)
            if t.ID == taskrunStatus.PipelineTaskName {
                task = &p.Stages[i].Tasks[j]
                break
            }
        }
        if task != nil {
            break
        }
    }
    if task == nil {
        return nil
    }
    task.TaskRunName = taskRunName
    task.StartTime = util.UnixMilli(status.StartTime)
    task.EndTime = util.UnixMilli(status.CompletionTime)
    task.Status = st
    task.SetActionsStatus(ctx, status)
    return task
}

func (p *PipelineRunPayload) NamespacedName() types.NamespacedName {
    return types.NamespacedName{
        Namespace: fmt.Sprintf("p6e-cx-%d", p.ProjectID),
        Name:      fmt.Sprintf("pipelinerun-%s", p.EventID),
    }
}
