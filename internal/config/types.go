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

import "net/url"

type Config struct {
	// Listen address binding the api
	Listen string `json:"listen"   yaml:"listen"`

	// Timeout is the time in minute until the watch for resource will be end
	Timeout int64 `json:"timeout"  yaml:"timeout"`

	// Auth is a list of credentials
	// it should be content the username and password
	Auth []Auth `json:"auth"     yaml:"auth"`

	// Webhooks is a list of subscriber where to send the status
	// it should be respect the url format
	// example for kafka :
	// kafka://localhost:9092?topic=MY-TOPIC
	Webhooks []Webhook `json:"webhooks" yaml:"webhooks"`
}

type Auth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Webhook struct {
	Type   string `json:"type" yaml:"type"`
	URLRaw string `json:"url"  yaml:"url"`
	url    *url.URL
}

var config = new(Config)

// IsAuth return a toggle against Auth config
func IsAuth() bool {
	return len(config.Auth) > 0
}
