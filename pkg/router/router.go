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
Created on 23/01/2021
*/

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/pkg/watch"
)

// New initializes the engine instance
func New() error {
	engine = gin.New()
	engine.Use(LogMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(CorrelationID())
	if config.IsAuth() {
		engine.Use(Auth())
	}
	AddPOST("/watch/play", watch.Play)
	return nil
}

// AddPOST adds handler and path to the engine
func AddPOST(relativePath string, handlers ...gin.HandlerFunc) {
	engine.POST(relativePath, handlers...)
}

// Run execute le gin router
func Run() error {
	return engine.Run(config.GetListen())
}
