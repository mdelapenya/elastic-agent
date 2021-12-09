// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cmd

import (
	"testing"
)

func TestAgent(t *testing.T) {
	// t.Run("test agent with subcommand", func(t *testing.T) {
	// 	streams, _, _, _ := cli.NewTestingIOStreams()
	// 	cmd := NewCommandWithArgs([]string{}, streams)
	// 	cmd.SetOutput(streams.Out)
	// 	cmd.Execute()
	// })

	// t.Run("test run subcommand", func(t *testing.T) {
	// 	streams, _, out, _ := cli.NewTestingIOStreams()
	// 	cmd := newRunCommandWithArgs(globalFlags{
	// 		PathConfigFile: filepath.Join("build", "elastic-agent.yml"),
	// 	}, []string{}, streams)
	// 	cmd.SetOutput(streams.Out)
	// 	cmd.Execute()
	// 	contents, err := ioutil.ReadAll(out)
	// 	if !assert.NoError(t, err) {
	// 		return
	// 	}
	// 	assert.True(t, strings.Contains(string(contents), "Hello I am running"))
	// })
}
