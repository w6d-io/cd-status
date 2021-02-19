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

import tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

const (
	pendingState = "---"
)

func GetSteps(status tkn.TaskRunStatus) (steps []Step) {
	log := logger.WithName("GetSteps")
	if (len(status.Steps) != 0 && status.Steps[0].Waiting == nil) &&
		(status.Steps[0].Terminated != nil || status.Steps[0].Running != nil) {
		log.V(2).Info("setting")
		for _, step := range status.Steps {
			log.V(2).Info("append", "name", step.Name, "status", StepReasonExists(step))
			steps = append(steps, Step{
				Name:   step.Name,
				Status: StepReasonExists(step),
			})
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
			return state.Terminated.Reason
		}
		return pendingState
	}
	return state.Waiting.Reason
}
