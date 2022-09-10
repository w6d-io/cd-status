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
	"context"
	"github.com/gin-gonic/gin"
	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/internal/tekton"
	"github.com/w6d-io/ci-status/pkg/handler/watch/play/pipelinerun"
	"github.com/w6d-io/ci-status/pkg/router"
	"github.com/w6d-io/x/logx"
	"io"
	"strings"
)

func init() {
	router.AddPOST("/watch/play", Play)
	AddWatcher(pipelinerun.KIND, pipelinerun.Scan)
}

// Play gets the play payload and determine the resource to scan
func Play(c *gin.Context) {
	log := logx.WithName(c.Request.Context(), "Play")
	var payload *tekton.PipelineRunPayload
	if err := c.BindJSON(&payload); err != nil && err != io.EOF {
		log.Error(err, "BindJSON")
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if payload.EventID == "" {
		c.JSON(400, gin.H{"status": "error", "message": "eventId is empty"})
		return
	}
	if payload.ProjectID == 0 {
		c.JSON(400, gin.H{"status": "error", "message": "projectId is empty"})
		return
	}
	go func(payload *tekton.PipelineRunPayload) {
		corId := c.Writer.Header().Get(config.CorrelationId)
		ip := GetClientIP(c)
		ctx := context.WithValue(context.Background(), logx.CorrelationID, corId)
		ctx = context.WithValue(ctx, logx.IPAddress, ip)
		ctx = context.WithValue(ctx, logx.Method, c.Request.Method)
		ctx = context.WithValue(ctx, logx.Kind, "pipelinerun")
		scan = scans[pipelinerun.KIND]
		err := scan(ctx, payload)
		if err != nil {
			logx.WithName(ctx, "Play").Error(err, "Scan resource")
			//c.JSON(403, gin.H{"status": "error", "message": "scan resource failed"})
			return
		}
	}(payload)
	c.JSON(200, gin.H{"status": "ok", "message": "scan launched"})
}

// AddWatcher inserts method to scans map
func AddWatcher(name string, f func(context.Context, *tekton.PipelineRunPayload) error) {
	scans[name] = f
}

// GetClientIP gets the correct IP for the end client instead of the proxy
func GetClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}
