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
Created on 05/02/2021
*/

package handler

import (
    "github.com/w6d-io/ci-status/pkg/handler/health"
    "github.com/w6d-io/ci-status/pkg/handler/watch"
)

func init() {
    _ = watch.Watch{}
    _ = health.Healthy{}
}

type Handler struct{}
