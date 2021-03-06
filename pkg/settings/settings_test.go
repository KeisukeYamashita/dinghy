/*
* Copyright 2019 Armory, Inc.

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

package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureSettings(t *testing.T) {
	cases := map[string]struct {
		defaults     Settings
		overrides    Settings
		expected     func(*testing.T, *Settings)
		expectsError bool
	}{
		"happy path": {
			defaults: Settings{
				GitHubToken: "12345",
			},
			overrides: Settings{
				GitHubToken: "45678",
			},
			expected: func(t *testing.T, settings *Settings) {
				assert.Equal(t, settings, &Settings{
					ParserFormat: "json",
					GitHubToken:  "45678",
				})
			},
		},
		"defaults no overrides": {
			defaults:  NewDefaultSettings(),
			overrides: Settings{},
			expected: func(t *testing.T, settings *Settings) {
				assert.NotEmpty(t, settings)
			},
		},
		"defaults with spinnaker settings overridden": {
			defaults: NewDefaultSettings(),
			overrides: Settings{
				spinnakerSupplied: spinnakerSupplied{
					Redis: Redis{
						BaseURL:  "12345:6789",
						Password: "",
					},
				},
			},
			expected: func(t *testing.T, settings *Settings) {
				assert.Equal(t, "12345:6789", settings.Redis.BaseURL)
				assert.Empty(t, settings.Redis.Password)
			},
		},
	}

	for testName, c := range cases {
		t.Run(testName, func(t *testing.T) {
			s, err := configureSettings(c.defaults, c.overrides)
			if !assert.Equal(t, c.expectsError, err != nil) {
				return
			}
			c.expected(t, s)
		})
	}
}

func TestDecodeProfilesToSettings(t *testing.T) {
	cases := map[string]struct {
		input        map[string]interface{}
		expected     Settings
		expectsError bool
	}{
		"happy path": {
			input: map[string]interface{}{
				"redis": map[string]interface{}{
					"baseUrl": "12345",
				},
			},
			expected: Settings{
				spinnakerSupplied: spinnakerSupplied{
					Redis: Redis{
						BaseURL: "12345",
					},
				},
			},
		},
	}

	for testName, c := range cases {
		t.Run(testName, func(t *testing.T) {
			var decoded Settings
			err := decodeProfilesToSettings(c.input, &decoded)
			if c.expectsError {
				assert.NotNil(t, err)
				return
			}
			assert.Equal(t, c.expected, decoded)

		})
	}
}

func TestSettings_GetRepoConfig(t *testing.T) {
	cases := map[string]struct {
		settings Settings
		provider string
		repo     string
		expected *RepoConfig
	}{
		"happy path": {
			settings: Settings{
				RepoConfig: []RepoConfig{
					{
						Provider: "github",
						Repo:     "ghrepo",
						Branch:   "ghbranch",
					},
					{
						Provider: "bitbucket",
						Repo:     "bbrepo",
						Branch:   "bbbranch",
					},
				},
			},
			provider: "bitbucket",
			repo:     "bbrepo",
			expected: &RepoConfig{
				Provider: "bitbucket",
				Repo:     "bbrepo",
				Branch:   "bbbranch",
			},
		},
		"not found": {
			settings: Settings{
				RepoConfig: []RepoConfig{
					{
						Provider: "github",
						Repo:     "ghrepo",
						Branch:   "ghbranch",
					},
					{
						Provider: "bitbucket",
						Repo:     "bbrepo",
						Branch:   "bbbranch",
					},
				},
			},
			provider: "stash",
			repo:     "repo",
			expected: nil,
		},
		"no repo configuration": {
			settings: Settings{RepoConfig: []RepoConfig{}},
			expected: nil,
		},
	}

	for testName, c := range cases {
		t.Run(testName, func(t *testing.T) {
			actual := c.settings.GetRepoConfig(c.provider, c.repo)
			assert.Equal(t, c.expected, actual)
		})
	}
}
