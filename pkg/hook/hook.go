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
Created on 07/02/2021
*/
package hook

import (
	"errors"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-status/internal/config"
	"net/url"
)

// Send loop into all the webhooks url. for each it get the function by the scheme and run the method/function associated
func Send(payload interface{}, logger logr.Logger) error {
	for _, wh := range config.GetWebhooks() {
		if wh.URL != nil {
			log := logger.WithValues("scheme", wh.URL.Scheme)
			f, ok := subscriber[wh.URL.Scheme]
			if !ok {
				log.Error(errors.New("scheme not supported"), "payload not sent")
			}
			if err := f(payload, wh.URL); err != nil {
				log.Error(err, "exec function failed")
			}
		}
	}
	return nil
}

// Subscribe adds the protocol Send function to the subscriber list
func Subscribe(name string, f func(interface{}, *url.URL) error) {
	subscriber[name] = f
}

