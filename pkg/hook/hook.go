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
	"fmt"
	"github.com/go-logr/logr"
	"github.com/w6d-io/ci-status/internal/config"
	"github.com/w6d-io/ci-status/pkg/hook/kafka"
	"net/url"
)

func init() {
	Subscribe("kafka", kafka.Send)
}

func Send(payload interface{}, logger logr.Logger) error {
	logger.V(1).Info("to send", "payload", payload)
	go func(payload interface{}, logger logr.Logger) {
		if err := DoSend(payload, logger); err != nil {
			logger.Error(err, "DoSend")
			return
		}
	}(payload, logger)
	return nil
}

// Send loop into all the webhooks url. for each it get the function by the scheme and run the method/function associated
func DoSend(payload interface{}, logger logr.Logger) error {
	log := logger.WithName("HookSend")
	whs := config.GetWebhooks()
	errc := make(chan error, len(whs))
	quit := make(chan struct{})
	defer close(quit)

	for _, wh := range whs {
		if wh.URL == nil {
			continue
		}
		log := log.WithValues("scheme", wh.URL.Scheme)
		go func(payload interface{}, URL *url.URL) {
			f, ok := subscriber[URL.Scheme]
			if !ok {
				err := fmt.Errorf("scheme %v not supported", URL.Scheme)
				log.Error(err, "payload not sent")
				return
			}
			logg := log.WithValues("url", URL)
			select {
			case errc <- f(payload, URL):
				logg.Info("sent")
			case <-quit:
				logg.Info("quit")
			}
		}(payload, wh.URL)
	}
	for range whs {
		if err := <-errc; err != nil {
			log.Error(err, "Sent failed")
			return err
		}
	}

	return nil
}

// Subscribe adds the protocol Send function to the subscriber list
func Subscribe(name string, f func(interface{}, *url.URL) error) {
	subscriber[name] = f
}
