/*
* Copyright 2020 Armory, Inc.

* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at

*    http://www.apache.org/licenses/LICENSE-2.0

* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package database

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

// RedisCacheReadOnly maintains a dependency graph inside Redis
type SQLReadOnly struct {
	Client *SQLClient
	Logger *log.Entry
	ctx    context.Context
	stop   chan os.Signal
}

// SetDeps sets dependencies for a parent
func (c *SQLReadOnly) SetDeps(parent string, deps []string) {

}

// GetRoots grabs roots
func (c *SQLReadOnly) GetRoots(url string) []string {
	return returnRoots(c.Client, url)
}

// Set RawData
func (c *SQLReadOnly) SetRawData(url string, rawData string) error {
	return nil
}

func (c *SQLReadOnly) GetRawData(url string) (string, error) {
	return returnRawData(c.Client, url)
}

// Clear clears everything
func (c *SQLReadOnly) Clear() {
}
