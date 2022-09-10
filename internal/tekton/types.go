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
	"github.com/tektoncd/cli/pkg/cli"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/types"
)

// Field is key for logs in context
type Field int

const (
//Timeout Field = iota + 1
)

type Tekton struct {
	ProjectID   int64               `json:"project_id"`
	EventID     int64               `json:"eventId"`
	PipelineRun *PipelineRunPayload `json:"pipeline_run"`
}

type Data map[string]string

type Trigger struct {
	ID          string `json:"id"          bson:"id"          mapstructure:"id"`
	Ref         string `json:"ref"         bson:"ref"         mapstructure:"ref"`
	Type        string `json:"type"        bson:"type"        mapstructure:"type"`
	ComponentId string `json:"componentId" bson:"componentId" mapstructure:"componentId"`
	Data        Data   `json:"data"        bson:"data"        mapstructure:"data"`
}
type ConditionType string

type Condition struct {
	Id   string        `json:"id"   bson:"id"`
	Ref  string        `json:"ref"  bson:"ref"`
	Type ConditionType `json:"type" bson:"type"`
	When string        `json:"when" bson:"when"`
}

type Conditions [][]Condition

type Action struct {
	ID           string                 `json:"id"           bson:"id"           mapstructure:"id"`
	Name         string                 `json:"name"         bson:"name"         mapstructure:"name"`
	ComponentID  string                 `json:"componentId"  bson:"componentId"  mapstructure:"componentId"`
	Ref          string                 `json:"ref"          bson:"ref"          mapstructure:"ref"`
	Data         map[string]string      `json:"data"         bson:"data"         mapstructure:"data"`
	Params       map[string]interface{} `json:"params"       bson:"params"       mapstructure:"params"`
	Environments map[string]string      `json:"environments" bson:"environments" mapstructure:"environments"`
	Status       string                 `json:"status"       bson:"status"       mapstructure:"status"`
	Reason       string                 `json:"reason"       bson:"reason"       mapstructure:"reason"`
	StartTime    int64                  `json:"startTime"    bson:"startTime"    mapstructure:"startTime"`
	EndTime      int64                  `json:"endTime"      bson:"endTime"      mapstructure:"endTime"`
}

type Task struct {
	TaskRunName   string     `json:"-"             bson:"-"             mapstructure:"-"`
	ID            string     `json:"id"            bson:"id"            mapstructure:"id"`
	Name          string     `json:"name"          bson:"name"          mapstructure:"name"`
	SkipOnFailure bool       `json:"skipOnFailure" bson:"skipOnFailure" mapstructure:"skipOnFailure"`
	Conditions    Conditions `json:"conditions"    bson:"conditions"    mapstructure:"conditions"`
	Actions       []Action   `json:"actions"       bson:"actions"       mapstructure:"actions"`
	StartTime     int64      `json:"startTime"     bson:"startTime"     mapstructure:"startTime"`
	EndTime       int64      `json:"endTime"       bson:"endTime"       mapstructure:"endTime"`
	Status        string     `json:"status"        bson:"status"        mapstructure:"status"`
	Reason        string     `json:"reason"        bson:"reason"        mapstructure:"reason"`
}

type Stage struct {
	ID        string `json:"id"         bson:"id"         mapstructure:"id"`
	Name      string `json:"name"       bson:"name"       mapstructure:"name"`
	Tasks     []Task `json:"tasks"      bson:"tasks"      mapstructure:"tasks"`
	Status    string `json:"status"     bson:"status"     mapstructure:"status"`
	EndTime   int64  `json:"endTime"    bson:"endTime"    mapstructure:"endTime"`
	StartTime int64  `json:"startTime"  bson:"startTime"  mapstructure:"startTime"`
}

type Commit struct {
	ID      string `json:"id"      bson:"id"      mapstructure:"id"`
	Ref     string `json:"ref"     bson:"ref"     mapstructure:"ref"`
	Message string `json:"message" bson:"message" mapstructure:"message"`
}

type PipelineRunPayload struct {
	ID               string    `json:"id"                  bson:"id"                mapstructure:"id"`
	Type             string    `json:"type"                bson:"type"              mapstructure:"type"`
	PipelineIDNumber string    `json:"pipelineIdNumber"    bson:"pipelineIdNumber"  mapstructure:"pipelineIdNumber"`
	ProjectID        int64     `json:"projectId"           bson:"projectId"         mapstructure:"projectId"`
	Name             string    `json:"name"                bson:"name"              mapstructure:"name"`
	Triggers         []Trigger `json:"triggers"            bson:"triggers"          mapstructure:"triggers"`
	Stages           []Stage   `json:"stages"              bson:"stages"            mapstructure:"stages"`
	Status           string    `json:"status"              bson:"status"            mapstructure:"status"`
	Reason           string    `json:"reason"              bson:"reason"            mapstructure:"reason"`
	StartTime        int64     `json:"startTime"           bson:"startTime"         mapstructure:"startTime"`
	EndTime          int64     `json:"endTime"             bson:"endTime"           mapstructure:"endTime"`
	LogUri           string    `json:"logUri"              bson:"logUri"            mapstructure:"logUri"`
	Complete         bool      `json:"complete"            bson:"complete"          mapstructure:"complete"`
	Force            bool      `json:"force"               bson:"-"                 mapstructure:"-"`
	Artifacts        bool      `json:"artifacts"           bson:"artifacts"         mapstructure:"artifacts"`
	TriggerId        string    `json:"triggerId,omitempty" bson:"triggerId"         mapstructure:"triggerId"`
	Commit           Commit    `json:"commit"              bson:"commit"            mapstructure:"commit"`
	EventID          string    `json:"eventId"             bson:"eventId"           mapstructure:"eventId"`
}

//type PipelineRunPayload struct {
//    NamespacedName types.NamespacedName `json:"namespaced_name"`
//    Status         string               `json:"status"`
//    Message        string               `json:"message"`
//    RepoURL        string               `json:"repo_url"`
//    Branch         string               `json:"branch"`
//    CommitID       string               `json:"commit_id"`
//    CommitMsg      string               `json:"commit_msg"`
//    TaskRuns       []TaskRunPayload     `json:"task_runs"`
//    ProjectID      int64                `json:"project_id"`
//    PipelineID     int64                `json:"pipeline_id"`
//    StartTime      int64                `json:"start_time"`
//    CompletionTime int64                `json:"completion_time"`
//}

type TaskRunPayload struct {
	NamespacedName types.NamespacedName `json:"namespaced_name"`
	Tasks          []Task               `json:"tasks"`
}

type Step struct {
	Name    string `json:"name"`
	RawName string `json:"raw_name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Tasks []*Task

var (
	logger       = ctrl.Log.WithName("tekton")
	tektonParams cli.TektonParams
	cls          *cli.Clients
)

// Reason ...
var Reason = map[string]string{
	"default":                         "",
	"couldntgetpipeline":              "Internal error",
	"invalidworkspacebindings":        "Internal error : Workspace resource missing",
	"invalidserviceaccountmappings":   "Internal error : invalid access",
	"invalidpipelineresourcebindings": "Internal error : invalid pipeline resource",
	"hyperparathyroidism":             "Internal error : Pipeline generation failed",
	"couldntgettask":                  "Internal error : task missing",
	"couldntgetresource":              "Internal error : missing resource",
	"pipelineruntimeout":              "Pipeline timeout",
	"couldntgetcondition":             "Internal error : Condition missing",
	"parametermissing":                "Internal error : Pipeline",
	"pipelinevalidationfailed":        "Internal error : Pipeline generation failed",
	"pipelineinvalidgraph":            "Internal error : Pipeline invalid",
	"pipelineruncancelled":            "Pipeline cancelled",
	"pipelineruncouldntcancel":        "Pipeline couldn't cancel",
	"pending":                         "Creating steps",
}
