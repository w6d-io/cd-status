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
	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/x/logx"
)

const (
	pendingState = "---"
)

func (t *Task) SetActionsStatus(ctx context.Context, status *tkn.TaskRunStatus) {
	log := logx.WithName(ctx, "Task.SetActionsStatus")
	log.V(2).Info("start")
	if (len(status.Steps) != 0 && status.Steps[0].Waiting == nil) &&
		(status.Steps[0].Terminated != nil || status.Steps[0].Running != nil) {
		log.V(1).Info("setting")
		for _, step := range status.Steps {
			log.V(2).Info("in loop", "step", step.Name)
			for i := range t.Actions {
				log.V(2).Info("in loop 2", "action", t.Actions[i].ID)
				if fmt.Sprintf("step-%s", t.Actions[i].ID) == step.Name || t.Actions[i].ID == step.Name {
					t.Actions[i].Status = StepReasonExists(step)
					log.V(1).Info("action", "name", t.Actions[i].Name, "id", t.Actions[i].ID, "status", t.Actions[i].Status)
				}
			}
		}
	}
	return
}

// StepReasonExists ...
func StepReasonExists(state tkn.StepState) string {
	logger.V(1).Info("StepReasonExists")
	if state.Waiting == nil {
		if state.Running != nil {
			return "Running"
		}
		if state.Terminated != nil {
			if state.Terminated.Reason == "Completed" {
				return "Succeeded"
			}
			return state.Terminated.Reason
		}
		return pendingState
	}
	return state.Waiting.Reason
}

// FormattedStepName ...
//func FormattedStepName(name string) string {
//
//    for _, gen := range genStepList {
//        if strings.HasPrefix(name, gen) {
//            return gen
//        }
//    }
//    return name
//}
