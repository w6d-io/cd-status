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
	"errors"
	"github.com/w6d-io/ci-status/internal/tekton"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/pkg/handler/watch/play/pipelinerun"
	"github.com/w6d-io/ci-status/pkg/router"
)

var (
	logger = ctrl.Log.WithName("watch")
)

func init() {
	router.AddPOST("/watch/play", Play)
	AddWatcher(pipelinerun.KIND, pipelinerun.Scan)
}

// Play gets the play payload and determine the resource to scan
func Play(c *gin.Context) {
	log := logger.WithName("Play")
	var payload Payload
	if err := c.BindJSON(&payload); err != nil {
		log.Error(err, "BindJSON")
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if payload.Object == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request"})
		return
	}
	log.V(1).Info("received", "payload", payload)
	if payload.Object.NamespacedName.Name == "" || payload.Object.NamespacedName.Namespace == "" {
		c.JSON(401, gin.H{"status": "error", "message": "name or/and namespace missing"})
		return
	}
	log = log.WithValues("kind", payload.Object.Kind)
	if scan, ok = scans[payload.Object.Kind]; !ok {
		log.Error(errors.New(payload.Object.Kind+" does not supported"), "BindJSON")
		c.JSON(402, gin.H{"status": "error", "message": payload.Object.Kind + " does not supported"})
		return
	}

	go func(kind string, nn types.NamespacedName, projectID, eventID int64, payload *tekton.PipelineRunPayload) {
		corId := c.Writer.Header().Get(config.CorrelationId)
		scanLog := ctrl.Log.WithValues("correlation_id", corId, "kind", kind)
		err := scan(scanLog, nn, projectID, eventID, payload)
		if err != nil {
			scanLog.Error(err, "Scan resource")
			//c.JSON(403, gin.H{"status": "error", "message": "scan resource failed"})
			return
		}
	}(payload.Object.Kind, payload.Object.NamespacedName, payload.ProjectID, payload.EventID, payload.Payload)
	c.JSON(200, gin.H{"status": "ok", "message": "scan launched"})
}

// AddWatcher inserts method to scans map
func AddWatcher(name string, f func(logr.Logger, types.NamespacedName, int64, int64, *tekton.PipelineRunPayload) error) {
	scans[name] = f
}
