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

package util

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// RemoveIndex takes a slice or an array and remove the element designated by the index
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// UnixMilli returns Unix time in millisecond
func UnixMilli(t *metav1.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1000000
}
