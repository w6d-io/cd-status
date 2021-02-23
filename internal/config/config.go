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
Created on 22/01/2021
*/

package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	WHSupport = []string{"kafka", "http", "https"}
)

func New(filename string) error {
	// TODO add dynamic configuration feature
	log := ctrl.Log.WithName("controllers").WithName("Config")
	log.V(1).Info("read config file")
	config = new(Config)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err, "error reading the configuration")
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Error(err, "Error unmarshal the configuration")
		return err
	}
	for i, wh := range config.Webhooks {
		if wh.URLRaw != "" {
			config.Webhooks[i].URL, err = url.Parse(wh.URLRaw)
			if err != nil {
				log.Error(err, "url "+wh.URLRaw+" parse failed")
				return err
			}
			if !IsInArray(config.Webhooks[i].URL.Scheme, WHSupport) {
				return fmt.Errorf("scheme %v not supported", config.Webhooks[i].URL.Scheme)
			}
			if config.Webhooks[i].URL.Scheme == "kafka" {
				values := config.Webhooks[i].URL.Query()
				if _, ok := values["topic"]; !ok {
					log.Error(errors.New("missing topic"), wh.URLRaw)
					return errors.New("missing topic")
				}
			}
		}
	}
	if config.Timeout == 0 {
		config.Timeout = int64(time.Minute * 60)
	}
	if config.Listen == "" {
		config.Listen = ":8080"
	}
	return nil
}

// GetConfig return the Config instance
func GetConfig() *Config {
	return config
}

// GetListen return the listening address for api
func GetListen() string {
	return config.Listen
}

// SetListen return the listening address for api
func SetListen(address string) {
	config.Listen = address
}

// GetTimeout return the timeout set in config
func GetTimeout() time.Duration {
	return time.Duration(config.Timeout) * time.Minute
}

// SetTimeout record the timeout in config
func SetTimeout(timeout int64) error {
	config.Timeout = timeout
	return nil
}

// SetAuth return the auth list from config
func SetAuth(auth []Auth) {
	config.Auth = auth
}

// GetAuth return the auth list from config
func GetAuth() []Auth {
	return config.Auth
}

// GetWebhook returns the list of webhook
func GetWebhooks() []Webhook {
	return config.Webhooks
}

// IsInArray checks if the needle in part of haystack
func IsInArray(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if elem == needle {
			return true
		}
	}
	return false
}
