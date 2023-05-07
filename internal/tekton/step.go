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
    tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
    "github.com/w6d-io/x/logx"
    "strings"
)

const (
    pendingState = "---"
)

var genStepList = []string{
    "create-dir-builddocker",
    "git-source",
    "image-digest-exporter",
}

func GetSteps(ctx context.Context, status tkn.TaskRunStatus) (steps []Step) {
    log := logx.WithName(ctx, "GetSteps")
    if (len(status.Steps) != 0 && status.Steps[0].Waiting == nil) &&
        (status.Steps[0].Terminated != nil || status.Steps[0].Running != nil) {
        log.V(1).Info("setting")
        for _, step := range status.Steps {
            log.V(1).Info("append", "name", step.Name, "status", StepReasonExists(ctx, step))
            steps = append(steps, Step{
                Name:    FormattedStepName(step.Name),
                RawName: step.Name,
                Status:  StepReasonExists(ctx, step),
            })
        }
    }
    return
}

// StepReasonExists ...
func StepReasonExists(ctx context.Context, state tkn.StepState) string {
    logx.WithName(ctx, "StepReasonExists").V(1).Info("call")
    if state.Waiting == nil {
        if state.Running != nil {
            return "Running"
        }
        if state.Terminated != nil {
            return state.Terminated.Reason
        }
        return pendingState
    }
    return state.Waiting.Reason
}

// FormattedStepName ...
func FormattedStepName(name string) string {

    for _, gen := range genStepList {
        if strings.HasPrefix(name, gen) {
            return gen
        }
    }
    return name
}
