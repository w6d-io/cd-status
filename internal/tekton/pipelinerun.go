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
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"github.com/tektoncd/cli/pkg/cli"
	"github.com/w6d-io/ci-status/internal/config"
	"k8s.io/apimachinery/pkg/watch"
)

// Supervise watches all pod event created by pipelinerun
func (t *Tekton) Supervise() error {
	log := logger.WithName("Supervise").WithValues("object", t.Namespaced.String())
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	timeout := time.NewTimer(time.Duration(config.GetTimeout()))

	for {
		select {
		case <-timeout.C:
			log.Info("timeout")
			return nil
		}
	}

}

// GetWatch gets the watch interface
func (t *Tekton) GetWatch() (w watch.Interface) {
	log := logger.WithName("GetWatch").WithValues("object", t.Namespaced.String())

	tknParam := cli.TektonParams{}
	tknParam.SetNamespace(t.Namespaced.Namespace)
	cs, err := tknParam.Clients()
	if err != nil {
		log.Error(err, "create tekton k8s api client")
		return nil
	}
	timeout := config.GetTimeout()
	opts := v1.ListOptions{
		FieldSelector:  fmt.Sprintf("metadata.name=%s", prName),
		TimeoutSeconds: &timeout,
	}
	switch expr {
	
	}
	return
}