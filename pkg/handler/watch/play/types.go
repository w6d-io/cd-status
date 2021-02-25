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
package play

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

//ScanFct func() error
//type ScanFct func(Payload, logr.Logger) error

type Payload struct {
	Object     *Object `json:"object,omitempty"`
	ProjectID  int64   `json:"project_id,omitempty"`
	PipelineID int64   `json:"pipeline_id,omitempty"`
	RepoURL    string  `json:"repo_url,omitempty"`
	Commit     Commit  `json:"ref,omitempty"`
	Stack      Stack   `json:"stack,omitempty"`
	Status     string  `json:"status,omitempty"`
}

// Commit contains all git information
type Commit struct {
	// SHA contains git commit SHA
	SHA string `json:"sha,omitempty"`

	// Ref contains git commit reference
	Ref string `json:"ref,omitempty"`

	// Message contains commit message
	Message string `json:"message,omitempty"`
}

type Object struct {
	Kind           string               `json:"kind,omitempty"`
	NamespacedName types.NamespacedName `json:"namespaced_name,omitempty"`
}

// Stack contains the language and package of the source
type Stack struct {
	// Language contains the repository language
	Language string `json:"language,omitempty"`

	// Package contains the package use in application
	Package string `json:"package,omitempty"`
}

type Interface interface {
	Scan()
}

var (
	scans   = make(map[string]func(logr.Logger, types.NamespacedName, int64, int64, string, string, string, string) error)
	scan    func(logr.Logger, types.NamespacedName, int64, int64, string, string, string, string) error
	payload Payload
	ok      bool
)
