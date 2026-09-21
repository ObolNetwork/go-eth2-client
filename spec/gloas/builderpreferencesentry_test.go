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

package gloas_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func validBuilderPreferencesEntry() *gloas.BuilderPreferencesEntry {
	return &gloas.BuilderPreferencesEntry{
		ProposerPubkey: phase0.BLSPubKey{0x0a},
		URL:            []byte("https://builder.example"),
		Auth: &gloas.SignedBuilderRequestAuth{
			Message:   &gloas.BuilderRequestAuth{Data: []byte("https://builder.example"), Slot: 123},
			Signature: phase0.BLSSignature{0x02},
		},
		MaxExecutionPayment: 4,
	}
}

func TestBuilderPreferencesEntryJSONRoundTrip(t *testing.T) {
	entry := validBuilderPreferencesEntry()

	data, err := json.Marshal(entry)
	require.NoError(t, err)

	var decoded gloas.BuilderPreferencesEntry

	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, entry, &decoded)
}

func TestBuilderPreferencesEntryJSONInvalid(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "missing proposer pubkey prefix",
			input:  `{"proposer_pubkey":"0a","url":"https://builder.example","auth":{"message":{"data":"0x01","slot":"1"},"signature":"0x` + zeroHex(96) + `"},"max_execution_payment":"1"}`,
			errMsg: "proposer public key missing 0x prefix",
		},
		{
			name:   "missing url",
			input:  `{"proposer_pubkey":"0x` + zeroHex(48) + `","url":"","auth":{"message":{"data":"0x01","slot":"1"},"signature":"0x` + zeroHex(96) + `"},"max_execution_payment":"1"}`,
			errMsg: "builder URL missing",
		},
		{
			name:   "missing auth",
			input:  `{"proposer_pubkey":"0x` + zeroHex(48) + `","url":"https://builder.example","max_execution_payment":"1"}`,
			errMsg: "builder authorization missing",
		},
		{
			name:   "invalid max execution payment",
			input:  `{"proposer_pubkey":"0x` + zeroHex(48) + `","url":"https://builder.example","auth":{"message":{"data":"0x01","slot":"1"},"signature":"0x` + zeroHex(96) + `"},"max_execution_payment":"nan"}`,
			errMsg: "invalid max execution payment",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var decoded gloas.BuilderPreferencesEntry

			err := json.Unmarshal([]byte(test.input), &decoded)
			require.ErrorContains(t, err, test.errMsg)
		})
	}
}

func TestBuilderPreferencesEntrySSZRoundTrip(t *testing.T) {
	entry := validBuilderPreferencesEntry()

	data, err := entry.MarshalSSZ()
	require.NoError(t, err)

	var decoded gloas.BuilderPreferencesEntry

	require.NoError(t, decoded.UnmarshalSSZ(data))
	require.Equal(t, entry, &decoded)

	root, err := entry.HashTreeRoot()
	require.NoError(t, err)
	require.NotEqual(t, [32]byte{}, root)
}

func TestBuilderPreferencesEntryYAMLRoundTrip(t *testing.T) {
	entry := validBuilderPreferencesEntry()

	data, err := entry.MarshalYAML()
	require.NoError(t, err)

	var decoded gloas.BuilderPreferencesEntry

	require.NoError(t, decoded.UnmarshalYAML(data))
	require.Equal(t, entry, &decoded)
}

func zeroHex(n int) string {
	b := make([]byte, n*2)
	for i := range b {
		b[i] = '0'
	}

	return string(b)
}
