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

package paths

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	isWindows := runtime.GOOS == "windows"
	testCases := []struct {
		Name        string
		Expected    string
		Actual      string
		ShouldMatch bool
	}{
		{"different paths", "/var/path/a", "/var/path/b", false},
		{"strictly same paths", "/var/path/a", "/var/path/a", true},
		{"strictly same win paths", `C:\Program Files\Elastic\Agent`, `C:\Program Files\Elastic\Agent`, true},
		{"case insensitive win paths", `C:\Program Files\Elastic\Agent`, `c:\Program Files\Elastic\Agent`, isWindows},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.ShouldMatch, ArePathsEqual(tc.Expected, tc.Actual))
		})
	}
}
