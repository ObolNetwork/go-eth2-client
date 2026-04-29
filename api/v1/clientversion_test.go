// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestClientVersionJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.clientVersionJSON",
		},
		{
			name:  "CodeMissing",
			input: []byte(`{"name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"}`),
			err:   "code missing",
		},
		{
			name:  "CodeWrongType",
			input: []byte(`{"code":true,"name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field clientVersionJSON.code of type string",
		},
		{
			name:  "NameMissing",
			input: []byte(`{"code":"LH","version":"v8.0.1","commit":"0xced49dd2"}`),
			err:   "name missing",
		},
		{
			name:  "NameWrongType",
			input: []byte(`{"code":"LH","name":true,"version":"v8.0.1","commit":"0xced49dd2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field clientVersionJSON.name of type string",
		},
		{
			name:  "VersionMissing",
			input: []byte(`{"code":"LH","name":"Lighthouse","commit":"0xced49dd2"}`),
			err:   "version missing",
		},
		{
			name:  "VersionWrongType",
			input: []byte(`{"code":"LH","name":"Lighthouse","version":true,"commit":"0xced49dd2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field clientVersionJSON.version of type string",
		},
		{
			name:  "CommitMissing",
			input: []byte(`{"code":"LH","name":"Lighthouse","version":"v8.0.1"}`),
			err:   "commit missing",
		},
		{
			name:  "CommitWrongType",
			input: []byte(`{"code":"LH","name":"Lighthouse","version":"v8.0.1","commit":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field clientVersionJSON.commit of type string",
		},
		{
			name:  "Good",
			input: []byte(`{"code":"LH","name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ClientVersion
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
