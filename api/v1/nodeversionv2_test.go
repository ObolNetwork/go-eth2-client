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

func TestNodeVersionV2JSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.nodeVersionV2JSON",
		},
		{
			name:  "BeaconNodeMissing",
			input: []byte(`{}`),
			err:   "beacon_node missing",
		},
		{
			name:  "BeaconNodeWrongType",
			input: []byte(`{"beacon_node":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type v1.clientVersionJSON",
		},
		{
			name:  "BeaconNodeInvalid",
			input: []byte(`{"beacon_node":{"name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"}}`),
			err:   "invalid JSON: code missing",
		},
		{
			name:  "ExecutionClientInvalid",
			input: []byte(`{"beacon_node":{"code":"LH","name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"},"execution_client":{"code":"NM","name":"Nethermind","commit":"0xc066aee2"}}`),
			err:   "invalid JSON: version missing",
		},
		{
			name:  "BeaconNodeOnly",
			input: []byte(`{"beacon_node":{"code":"LH","name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"}}`),
		},
		{
			name:  "Full",
			input: []byte(`{"beacon_node":{"code":"LH","name":"Lighthouse","version":"v8.0.1","commit":"0xced49dd2"},"execution_client":{"code":"NM","name":"Nethermind","version":"v1.35.8","commit":"0xc066aee2"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.NodeVersionV2
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
