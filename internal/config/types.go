/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 22/01/2021
*/

package config

type Config struct {
	// TODO Add auth part
}

var config = new(Config)

// IsAuth return a toggle against Auth config
func IsAuth() bool {
	// TODO Check auth entry in Config
	return false
}
