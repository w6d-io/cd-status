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
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gin-gonic/gin"
)

// engine is the gin Engine instance in router package
var (
	server = &http.Server{}
	engine = gin.New()
	logger = ctrl.Log.WithName("Router")
)
