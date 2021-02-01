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
package watch

import (
	"errors"
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gin-gonic/gin"
)

var (
	logger = ctrl.Log.WithName("watch")
)

// Play gets the play payload and determine the resource to scan
func Play(c *gin.Context) {
	log := logger.WithName("Play")
	var (
		payload Payload
		scan    func(Payload) error
		ok      bool
	)
	if err := c.BindJSON(&payload); err != nil {
		log.Error(err, "BindJSON")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	log = log.WithValues("kind", payload.Object.Kind)
	if scan, ok = scans[payload.Object.Kind]; !ok {
		log.Error(errors.New(payload.Object.Kind+" does not supported"), "BindJSON")
		c.JSON(http.StatusNotImplemented, gin.H{"status": "error", "message": payload.Object.Kind + " does not supported"})
		return
	}
	if err := scan(payload); err != nil {
		log.Error(err, "Scan resource")
		return
	}
}

// AddWatcher inserts method to scans map
func AddWatcher(name string, w func(payload Payload) error) {
	scans[name] = w
}
