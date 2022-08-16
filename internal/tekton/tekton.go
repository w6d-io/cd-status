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
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tektoncd/cli/pkg/cli"
	"github.com/w6d-io/ci-status/internal/config"
	"k8s.io/apimachinery/pkg/watch"
	"knative.dev/pkg/apis/duck/v1beta1"
)

// GetWatch gets the watch interface for pipelineruns, taskruns and pods
func (t *Tekton) GetWatch(kind string, name string) (w watch.Interface) {
	log := logger.WithName("GetWatch").WithValues("object", t.PipelineRun.NamespacedName.String(), "kind", kind)
	namespace := t.PipelineRun.NamespacedName.Namespace
	//tknParam := cli.TektonParams{}
	var err error
	timeout := int64(config.GetTimeout().Seconds())
	opts := metav1.ListOptions{
		FieldSelector:  fmt.Sprintf("metadata.name=%s", name),
		TimeoutSeconds: &timeout,
	}
	cs := t.GetClient(namespace)
	if cs == nil {
		log.Error(errors.New("get client return nil"), "GetClient")
		return
	}
	if cs.Tekton == nil {
		log.Error(errors.New("tekton client is nil"), "check tekton client")
		return nil
	}
	switch kind {
	case "pipelinerun", "pipelineruns":
		log.V(1).Info("get pipelinerun", "name", name)

		_, err = cs.Tekton.TektonV1beta1().PipelineRuns(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			log.V(1).Info("not found")
			return nil
		}
		if err != nil {
			log.Error(err, "get resource failed")
			return nil
		}
		w, err = cs.Tekton.TektonV1beta1().PipelineRuns(namespace).Watch(context.TODO(), opts)
	case "taskrun", "taskruns":
		log.V(1).Info("get taskrun", "name", name)
		_, err := cs.Tekton.TektonV1beta1().TaskRuns(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			log.V(1).Info("not found")
			return nil
		}
		if err != nil {
			log.Error(err, "get resource failed")
			return nil
		}
		w, err = cs.Tekton.TektonV1beta1().TaskRuns(namespace).Watch(context.TODO(), opts)
	case "po", "pod", "pods":
		log.V(1).Info("get pod", "name", name)
		_, err = cs.Kube.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			log.V(1).Info("not found")
			return nil
		}
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

// Len is a method for Sort
func (t Tasks) Len() int {
	return len(t)
}

// Swap is a method for Sort
func (t Tasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less is a sort method
func (t Tasks) Less(i, j int) bool {
	return t[i].StartTime < t[j].StartTime
}

// Condition returns a human readable text based on the status of the Condition
func Condition(c v1beta1.Conditions) (status string, reason string) {
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

// IsTerminated returns if a pod is terminated
func IsTerminated(c v1beta1.Conditions) bool {
	if len(c) == 0 {
		return false
	}
	switch c[0].Status {
	case corev1.ConditionFalse, corev1.ConditionTrue:
		return true
	default:
		return false
	}
}

func (t *Tekton) SetParam(params cli.TektonParams) {
	tektonParams = params
}

func (t *Tekton) SetClient(clients *cli.Clients) {
	cls = clients
}

func (t *Tekton) GetClient(namespace string) *cli.Clients {
	if cls != nil {
		return cls
	}
	var err error
	tektonParams = cli.TektonParams{}
	tektonParams.SetNamespace(namespace)
	cls, err = tektonParams.Clients()
	if err != nil {
		cls = nil
	}
	return cls
}
