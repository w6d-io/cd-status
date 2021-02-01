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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/types"
)

type Tekton struct {
	ProjectID  int64                `json:"project_id"`
	PipelineID int64                `json:"pipeline_id"`
	Namespaced types.NamespacedName `json:"namespaced"`
	TaskRun    TaskRun              `json:"task_run"`
}

type TaskRun struct {
	Name  string `json:"name"`
	Tasks Tasks  `json:"tasks"`
}

type Tasks []Task

type Task struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	Duration  string `json:"duration"`
	Message   string `json:"message"`
	Steps     []Step `json:"steps"`

	CompletionTimeRaw *metav1.Time `json:"completion_time_raw"`
	StartTimeRaw      *metav1.Time `json:"start_time_raw"`
}

type Step struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

var (
	logger = ctrl.Log.WithName("tekton")
)

type PipelineRunPayload struct {
	NamespacedName types.NamespacedName `json:"namespaced_name"`
	Status         string               `json:"status"`
	Message        string               `json:"message"`
}
