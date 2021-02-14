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
Created on 27/01/2021
*/
package tekton

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jonboulle/clockwork"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/formatted"
	"github.com/w6d-io/ci-status/internal/config"
	"k8s.io/apimachinery/pkg/watch"
	"knative.dev/pkg/apis/duck/v1beta1"
)

// GetWatch gets the watch interface for pipelineruns, taskruns and pods
func (t *Tekton) GetWatch(kind string) (w watch.Interface) {
	log := logger.WithName("GetWatch").WithValues("object", t.PipelineRun.NamespacedName.String(), "kind", kind)
	namespace := t.PipelineRun.NamespacedName.Namespace
	tknParam := cli.TektonParams{}
	tknParam.SetNamespace(namespace)
	cs, err := tknParam.Clients()
	if err != nil {
		log.Error(err, "create tekton k8s api client")
		return nil
	}
	timeout := int64(config.GetTimeout().Seconds())
	opts := metav1.ListOptions{
		FieldSelector:  fmt.Sprintf("metadata.name=%s", t.PipelineRun.NamespacedName.Name),
		TimeoutSeconds: &timeout,
	}
	switch kind {
	case "pipelinerun", "pipelineruns":
		w, err = cs.Tekton.TektonV1beta1().PipelineRuns(namespace).Watch(context.TODO(), opts)
	case "taskrun", "taskruns":
		w, err = cs.Tekton.TektonV1beta1().TaskRuns(namespace).Watch(context.TODO(), opts)
	case "po", "pod", "pods":
		w, err = cs.Kube.CoreV1().Pods(namespace).Watch(context.TODO(), opts)
	default:
		log.Error(errors.New("kind not supported"), "")
		return nil
	}
	if err != nil {
		log.Error(err, "get watch instance failed")
		return nil
	}
	return
}

func (t *Tekton) GetTask(pr *tkn.PipelineRun) (ts Tasks) {
	log := logger.WithName("GetTask").WithValues("object", t.PipelineRun.NamespacedName.String())
	if len(pr.Status.TaskRuns) > 0 {
		for _, taskrunStatus := range pr.Status.TaskRuns {
			task := Task{
				Name: taskrunStatus.PipelineTaskName,
				StartTime: formatted.Age(taskrunStatus.Status.StartTime,
					clockwork.NewRealClock()),
				Duration: formatted.Duration(taskrunStatus.Status.StartTime,
					taskrunStatus.Status.CompletionTime),
				Status:            formatted.Condition(taskrunStatus.Status.Conditions),
				StartTimeRaw:      taskrunStatus.Status.StartTime,
				CompletionTimeRaw: taskrunStatus.Status.CompletionTime,
			}
			ts = append(ts, task)
		}
		sort.Sort(ts)
		log.V(1).WithValues("elements", len(ts)).Info("tasks")
	}
	return ts
}

// Len is a method for Sort
func (t Tasks) Len() int {
	return len(t)
}

// Swap is a method for Sort
func (t Tasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t Tasks) Less(i, j int) bool {
	if t[j].StartTimeRaw == nil {
		return false
	}
	if t[i].StartTimeRaw == nil {
		return true
	}
	return t[j].StartTimeRaw.Before(t[i].StartTimeRaw)
}

// GetStatusReason returns a human readable text based on the status of the Condition
func GetStatusReason(c v1beta1.Conditions) (status string, reason string) {
	if len(c) == 0 {
		return "---", ""
	}

	switch c[0].Status {
	case corev1.ConditionFalse:
		status = "Failed"
	case corev1.ConditionTrue:
		status = "Succeeded"
	case corev1.ConditionUnknown:
		status = "Running"
	}
	cStatus := status

	if c[0].Reason != "" && c[0].Reason != status {

		if c[0].Reason == "PipelineRunCancelled" || c[0].Reason == "TaskRunCancelled" {
			reason = Reason["default"]
			if _, ok := Reason[strings.ToLower(c[0].Reason)]; ok {
				reason = Reason[strings.ToLower(c[0].Reason)]
			}
			status = "Cancelled"
		} else if c[0].Reason != status {
			reason = Reason["default"]
			if _, ok := Reason[strings.ToLower(c[0].Reason)]; ok {
				reason = Reason[strings.ToLower(c[0].Reason)]
			}
			status = cStatus
		}
	} else {
		status = cStatus
	}
	return
}